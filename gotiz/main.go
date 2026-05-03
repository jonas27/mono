// Command gotiz downloads cycling videos from tiz-cycling.tv and Discovery+.
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jmanser/gotiz/internal/downloader"
	"github.com/jmanser/gotiz/internal/dplus"
	"github.com/jmanser/gotiz/internal/scraper"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		flagOutput      = flag.String("output", ".", "download directory")
		flagDryRun      = flag.Bool("dry-run", false, "print URLs without downloading")
		flagDplusConfig = flag.String("dplus-config", ".discovery.yaml", "Discovery+ config file")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gotiz [flags] <URL>...\n\n")
		fmt.Fprintf(os.Stderr, "Supported URLs:\n")
		fmt.Fprintf(os.Stderr, "  tiz-cycling.tv:    video pages (/video/...) or category pages (/categories/...)\n")
		fmt.Fprintf(os.Stderr, "  play.discoveryplus.com or www.discoveryplus.com: video URLs with UUID paths\n\n")
		fmt.Fprintf(os.Stderr, "Discovery+ uses a headless browser; set 'username' and 'password' in the config.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return fmt.Errorf("at least one URL is required")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Fprintln(os.Stderr, "\ncancelling...")
		cancel()
	}()

	if !*flagDryRun {
		if err := os.MkdirAll(*flagOutput, 0o750); err != nil {
			return fmt.Errorf("create output directory: %w", err)
		}
	}

	sc := scraper.New()
	dl := downloader.New()

	dplusCfg, _ := loadDiscoveryConfig(*flagDplusConfig)
	dplusClient := dplus.New(dplusCfg)

	for _, rawURL := range args {
		if isDiscoveryURL(rawURL) {
			if err := handleDiscovery(ctx, dplusClient, rawURL, *flagOutput, *flagDryRun); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				fmt.Fprintf(os.Stderr, "error: %s: %v\n", rawURL, err)
			}
			continue
		}

		videoURLs, err := sc.Resolve(ctx, rawURL)
		if err != nil {
			return fmt.Errorf("resolve %s: %w", rawURL, err)
		}

		for _, videoURL := range videoURLs {
			info, err := sc.FetchVideo(ctx, videoURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: skip %s: %v\n", videoURL, err)
				continue
			}

			if *flagDryRun {
				fmt.Printf("[dry-run] %s\n  %s\n", info.Title, info.DirectURL)
				continue
			}

			fmt.Printf("downloading: %s\n", info.Title)

			if err := dl.Download(ctx, info.DirectURL, *flagOutput, info.Filename); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				fmt.Fprintf(os.Stderr, "error: %s: %v\n", info.Title, err)
			}
		}
	}

	return nil
}

func handleDiscovery(ctx context.Context, client *dplus.Client, rawURL, dir string, dryRun bool) error {
	info, err := client.Resolve(ctx, rawURL)
	if err != nil {
		return err
	}

	if dryRun {
		drm := ""
		if info.DRMProtected {
			drm = " [DRM]"
		}
		fmt.Printf("[dry-run] %s\n  %s (%s%s)\n", info.Title, info.StreamURL, info.Kind, drm)
		return nil
	}

	return dplus.Download(ctx, info, dir)
}

func isDiscoveryURL(rawURL string) bool {
	return strings.Contains(rawURL, "discoveryplus.com")
}

// loadDiscoveryConfig reads a flat YAML file with key: value pairs.
func loadDiscoveryConfig(path string) (dplus.Config, error) {
	f, err := os.Open(path) //nolint:gosec
	if err != nil {
		return dplus.Config{}, err
	}
	defer f.Close()

	var cfg dplus.Config
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)

		switch k {
		case "username":
			cfg.Username = v
		case "password":
			cfg.Password = v
		case "token":
			cfg.Token = v
		}
	}

	return cfg, sc.Err()
}
