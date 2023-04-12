package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

const (
	// exitFail is the exit code if the program
	// fails.
	exitFail = 1
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	root := flag.String("root", "$HOME", "Root folder for kubernetes configs")
	flag.Parse()

	if err := run(*root); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run(dir string) error {
	dirs, err := walkDir(dir)
	if err != nil {
		return err
	}
	option, err := customAuth()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, v := range dirs {
		wg.Add(1)
		go func(v string) {
			defer wg.Done()
			yaml, err := kustomize(v)
			if err != nil {
				log.Printf("Can't build %s", v)
				// log.Println(err)
			}
			images := parseImages(yaml)
			err = checkImages(images, option)
			if err != nil {
				log.Printf("dir %s: %s", v, err)
			}

		}(v)
	}
	wg.Wait()
	return nil
}

func walkDir(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == "staging" {
				p := filepath.Join(path, "kustomization.yaml")
				if _, err := os.Stat(p); err == nil {
					files = append(files, path)
				}
			}
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
		return nil, err
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
func checkImages(images []string, option remote.Option) error {
	for _, img := range images {
		if strings.Contains(img, "@") || strings.Count(img, ":") == 1 {
			img_split := strings.Split(img, "@")
			desc, err := remoteImage(img_split[0], option)
			if err != nil {
				continue // image has either no digest or no version
			}
			if len(img_split) > 1 && img_split[1] != desc.Digest.String() {
				return fmt.Errorf("%s digest does not match version tag", img)
			}
		}
	}
	return nil
}

func remoteImage(img string, option remote.Option) (*remote.Descriptor, error) {
	ref, err := name.ParseReference(img)
	if err != nil {
		return nil, err
	}
	var desc *remote.Descriptor
	if strings.Contains(img, "custom.docker.com") {
		desc, err = remote.Get(ref, option)
	} else {
		desc, err = remote.Get(ref)
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
