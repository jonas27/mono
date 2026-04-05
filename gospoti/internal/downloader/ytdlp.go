package downloader

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// YTDLPRunner wraps the yt-dlp subprocess for downloading audio.
type YTDLPRunner struct {
	quality    string
	cookieFile string
	ytdlpPath  string
	ffmpegPath string
}

// NewYTDLPRunner creates a YTDLPRunner using the given binary paths and audio quality (kbps).
func NewYTDLPRunner(quality, ytdlpPath, ffmpegPath string) *YTDLPRunner {
	return &YTDLPRunner{
		quality:    quality,
		cookieFile: "youtube-cookies.txt",
		ytdlpPath:  ytdlpPath,
		ffmpegPath: ffmpegPath,
	}
}

// DownloadSearch downloads the first YouTube search result for the query string
// and saves it as an MP3 file at outPath.
func (y *YTDLPRunner) DownloadSearch(ctx context.Context, query, outPath string) error {
	return y.download(ctx, "ytsearch1:"+query, outPath)
}

// DownloadURL downloads the audio from a direct YouTube URL and saves it as MP3 at outPath.
func (y *YTDLPRunner) DownloadURL(ctx context.Context, videoURL, outPath string) error {
	return y.download(ctx, videoURL, outPath)
}

func (y *YTDLPRunner) download(ctx context.Context, source, outPath string) error {
	// Strip .mp3 extension – yt-dlp appends its own extension before ffmpeg post-processing.
	outTemplate := outPath
	if filepath.Ext(outTemplate) == ".mp3" {
		outTemplate = outTemplate[:len(outTemplate)-4]
	}

	args := []string{
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", y.quality + "K",
		"--output", outTemplate + ".%(ext)s",
		"--quiet",
		"--no-progress",
		"--ffmpeg-location", y.ffmpegPath,
	}

	if _, err := os.Stat(y.cookieFile); err == nil {
		args = append(args, "--cookies", y.cookieFile)
	}

	args = append(args, "--", source)

	cmd := exec.CommandContext(ctx, y.ytdlpPath, args...) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp download failed for %q: %w", source, err)
	}

	return nil
}
