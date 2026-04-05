// Command gospoti downloads music as MP3 files from Spotify playlists, albums, and tracks.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmanser/gospoti/internal/binaries"
	"github.com/jmanser/gospoti/internal/config"
	"github.com/jmanser/gospoti/internal/downloader"
	"github.com/jmanser/gospoti/internal/lang"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		flagOutput  = flag.String("output", "", "download folder (default: from config or 'downloads')")
		flagQuality = flag.String("quality", "", "128|192|256|320 (default: from config or '192')")
		flagLang    = flag.String("lang", "en", "en|bg|es")
		flagSetup   = flag.Bool("setup", false, "configure API keys interactively")
		flagArtist  = flag.String("artist", "", "override artist name for all downloads")
	)

	flag.Parse()

	cfg := config.Load()

	langCode := lang.Code(*flagLang)

	// Handle --setup flag.
	if *flagSetup {
		if err := setupKeys(&cfg); err != nil {
			return err
		}

		fmt.Println("API keys saved.")

		return nil
	}

	// Resolve output dir and quality.
	outDir := cfg.DownloadPath
	if *flagOutput != "" {
		outDir = *flagOutput
	}

	quality := cfg.Quality
	if *flagQuality != "" {
		quality = *flagQuality
	}

	// Resolve yt-dlp and ffmpeg binaries (bundled or from PATH).
	toolPaths, err := binaries.Resolve()
	if err != nil {
		return err
	}

	// Ensure output directory exists.
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return fmt.Errorf("create output directory %s: %w", outDir, err)
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return fmt.Errorf("Spotify URL argument is required")
	}

	rawURL := args[0]

	// Prompt for Spotify keys if missing.
	if cfg.APIKeys.SpotifyID == "" || cfg.APIKeys.SpotifySecret == "" {
		fmt.Println("Spotify API keys are required.")

		if err := setupKeys(&cfg); err != nil {
			return err
		}
	}

	dl, err := downloader.New(
		cfg.APIKeys.SpotifyID,
		cfg.APIKeys.SpotifySecret,
		cfg.APIKeys.GeniusToken,
		outDir,
		quality,
		*flagArtist,
		langCode,
		toolPaths.YTDlp,
		toolPaths.FFmpeg,
	)
	if err != nil {
		return fmt.Errorf("create downloader: %w", err)
	}

	// Set up graceful shutdown on Ctrl+C.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\n" + lang.Get("cancelling_message", langCode))
		cancel()
	}()

	return dl.DownloadSpotify(ctx, rawURL)
}

// setupKeys prompts the user for API keys and saves them to config.json.
func setupKeys(cfg *config.Config) error {
	keys, err := config.PromptForKeys(cfg.APIKeys)
	if err != nil {
		return fmt.Errorf("prompt for API keys: %w", err)
	}

	cfg.APIKeys = keys

	if err := config.Save(*cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}
