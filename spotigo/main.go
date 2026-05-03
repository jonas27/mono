// Command spotigo streams music directly from Spotify, writing raw PCM files.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmanser/spotigo/internal/creds"
	"github.com/jmanser/spotigo/internal/player"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		flagCreds  = flag.String("creds", ".creds.yaml", "credentials file")
		flagOutput = flag.String("output", "downloads", "output directory (f32le PCM, 44100 Hz, stereo)")
		flagAlbum  = flag.String("album", "", "override album/playlist name used as output subdirectory")
	)
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return fmt.Errorf("spotify URL or URI required")
	}

	c, err := creds.Load(*flagCreds)
	if err != nil {
		return err
	}

	sess, err := player.New(c, *flagCreds)
	if err != nil {
		return err
	}
	defer sess.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() { <-sig; cancel() }()

	return sess.Run(ctx, args[0], *flagOutput, *flagAlbum)
}
