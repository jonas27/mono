package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"golang.org/x/sync/errgroup"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

const (
	// exitFail is the exit code if the program fails.
	exitFail = 1
)

type imagesMutex struct {
	mu     sync.Mutex
	images []string
}

var Version = "development"

func main() {
	var root string

	rootCmd := &cobra.Command{
		Use:   "kustomize-check",
		Short: "Check all dirs under `root` if they are kustomize compilable",
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(root); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(exitFail)
			}
		},
	}

	rootCmd.Flags().StringVarP(&root, "root", "", "/home/joe/repos/k8s-setup", "Root folder for kubernetes configs")
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version %s\n", Version)
		},
	}

	rootCmd.AddCommand(versionCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitFail)
	}
}

func run(root string) error {
	images, err := runKustomize(root)
	if err != nil {
		return err
	}
	err = runImages(images)
	if err != nil {
		return err
	}
	return nil
}
func runKustomize(root string) (*imagesMutex, error) {
	dirs, err := walkDir(root)
	if err != nil {
		return &imagesMutex{}, err
	}

	ctx := context.Background()
	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(100)

	var images imagesMutex

	log.Printf("len dirs: %d", len(dirs))
	for _, v := range dirs {
		dir := v
		group.Go(func() error {
			if ctx.Err() == nil {
				yaml, err := kustomize(dir)
				if err != nil {
					return fmt.Errorf("Can't build %s", dir)
				}
				is := parseImages(yaml)
				images.mu.Lock()
				images.images = append(images.images, is...)
				images.mu.Unlock()
			}
			return nil
		})
	}
	if err = group.Wait(); err != nil {
		return &imagesMutex{}, err
	}
	log.Println("done with kustomize")
	return &images, nil
}

func runImages(images *imagesMutex) error {
	option, err := customAuth()
	if err != nil {
		return err
	}
	// get unique images
	imagesUnique := make(map[string]bool)
	for _, v := range images.images {
		imagesUnique[v] = true
	}
	keys := []string{}
	for k, _ := range imagesUnique {
		keys = append(keys, k)
	}
	log.Printf("len images: %d", len(keys))

	// check unique images
	ctx := context.Background()
	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(100)
	for _, k := range keys {
		kk := k
		group.Go(func() error {
			if ctx.Err() == nil {
				err := checkImage(kk, option)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return nil
}

func walkDir(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "kustomization.yaml" {
			files = append(files, filepath.Dir(path))
		}
		return nil
	})
	return files, err
}

func kustomize(dir string) (yaml []byte, err error) {
	tpc := types.PluginConfig{
		HelmConfig: types.HelmConfig{
			Enabled: true,
			Command: "helm",
		},
	}
	kOpts := krusty.Options{
		PluginConfig: &tpc,
	}
	k := krusty.MakeKustomizer(&kOpts)
	fSys := filesys.MakeFsOnDisk()
	m, err := k.Run(fSys, dir)
	if err != nil {
		return nil, fmt.Errorf("dir:%s\n%s", dir, err)
	}
	yaml, err = m.AsYaml()
	if err != nil {
		return nil, err
	}
	return yaml, nil
}

func parseImages(yml []byte) []string {
	yaamls := strings.Split(string(yml), "\n")
	images := []string{}
	for _, s := range yaamls {
		if strings.Contains(s, "image:") {
			ss := strings.Split(s, "image:")
			images = append(images, strings.TrimSpace(ss[1]))
		}
	}
	return images
}

// not really happy about how the image is calculated; via strings methods.
// At least there should be a check, maybe the go-container lib has a method to convert string to image url.
func checkImage(image string, option remote.Option) error {
	if strings.Contains(image, "@") || strings.Count(image, ":") == 1 {
		img_split := strings.Split(image, "@")
		desc, err := remoteImage(img_split[0], option)
		if err != nil {
			return err // image has either no digest or no version
		}
		if len(img_split) > 1 && img_split[1] != desc.Digest.String() {
			return fmt.Errorf("%s digest does not match version tag", image)
		}
	} else {
		_, err := remoteImage(image, option)
		if err != nil {
			return err
		}
	}
	return nil
}

func remoteImage(img string, option remote.Option) (*v1.Descriptor, error) {
	ref, err := name.ParseReference(img)
	if err != nil {
		return nil, err
	}
	var desc *v1.Descriptor
	if strings.Contains(img, "custom.docker.com") {
		desc, err = remote.Head(ref, option)
	} else {
		desc, err = remote.Head(ref)
	}
	if err != nil {
		return nil, err
	}
	return desc, nil
}

func customAuth() (remote.Option, error) {
	auth := authn.Basic{
		Username: "test",
		Password: "password",
	}
	conf, err := auth.Authorization()
	if err != nil {
		return nil, err
	}
	conf.Auth = "SOMETHING"
	option := remote.WithAuth(&auth)
	if err != nil {
		return nil, err
	}
	return option, nil
}
