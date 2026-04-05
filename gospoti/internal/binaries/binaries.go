// Package binaries resolves the paths to the yt-dlp and ffmpeg executables.
// When the binary is compiled with -tags bundled the tools are extracted from
// embedded assets; otherwise the PATH is searched.
package binaries

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Paths holds resolved paths to yt-dlp and ffmpeg.
type Paths struct {
	YTDlp  string
	FFmpeg string
}

// Resolve returns paths to usable yt-dlp and ffmpeg binaries.
// With -tags bundled the embedded binaries are extracted to the OS cache
// directory on first use and those paths are returned.
// Without the tag the PATH is searched.
func Resolve() (Paths, error) {
	ytdlp, err := resolveTool("yt-dlp", ytdlpBin)
	if err != nil {
		return Paths{}, err
	}

	ffmpeg, err := resolveTool("ffmpeg", ffmpegBin)
	if err != nil {
		return Paths{}, err
	}

	return Paths{YTDlp: ytdlp, FFmpeg: ffmpeg}, nil
}

func resolveTool(name string, embedded []byte) (string, error) {
	if len(embedded) > 0 {
		return extractBundled(name, embedded)
	}

	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("%s not found in PATH: %w", name, err)
	}

	return path, nil
}

func extractBundled(name string, data []byte) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}

	sum := sha256.Sum256(data)
	dir := filepath.Join(cacheDir, "gospoti", hex.EncodeToString(sum[:8]))

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create cache dir for %s: %w", name, err)
	}

	binName := name
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	dest := filepath.Join(dir, binName)

	if _, err := os.Stat(dest); err == nil {
		return dest, nil
	}

	if err := os.WriteFile(dest, data, 0o755); err != nil {
		return "", fmt.Errorf("extract %s: %w", name, err)
	}

	return dest, nil
}
