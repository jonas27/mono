// Package downloader fetches video files over HTTP with progress display and resume support.
package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Downloader fetches and saves video files from direct HTTP/HTTPS URLs.
type Downloader struct {
	client *http.Client
}

// New returns a Downloader with no timeout (suitable for large video files).
func New() *Downloader {
	return &Downloader{
		client: &http.Client{},
	}
}

// Download fetches src and writes it to filepath.Join(dir, filename).
// If a partial file exists it attempts to resume via HTTP Range.
// Progress is printed to stdout with carriage-return overwriting.
func (d *Downloader) Download(ctx context.Context, src, dir, filename string) error {
	dest := filepath.Join(dir, filename)

	var offset int64
	if fi, err := os.Stat(dest); err == nil {
		offset = fi.Size()
	}

	req, err := buildRequest(ctx, src, offset)
	if err != nil {
		return err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		offset = 0
	case http.StatusPartialContent:
		// continue from offset
	case http.StatusRequestedRangeNotSatisfiable:
		fmt.Printf("  already complete\n")
		return nil
	default:
		return fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	total := resp.ContentLength
	if total > 0 {
		total += offset
	}

	if offset > 0 && total > 0 && offset >= total {
		fmt.Printf("  already complete (%s)\n", fmtBytes(offset))
		return nil
	}

	f, err := openDest(dest, offset, resp.StatusCode)
	if err != nil {
		return err
	}
	defer f.Close()

	return copyWithProgress(f, resp.Body, offset, total)
}

func buildRequest(ctx context.Context, src string, offset int64) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; gotiz/1.0)")

	if offset > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}

	return req, nil
}

func openDest(dest string, offset int64, status int) (*os.File, error) {
	if offset > 0 && status == http.StatusPartialContent {
		f, err := os.OpenFile(dest, os.O_WRONLY|os.O_APPEND, 0o644) //nolint:mnd
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", dest, err)
		}

		return f, nil
	}

	f, err := os.Create(dest)
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", dest, err)
	}

	return f, nil
}

func copyWithProgress(dst io.Writer, src io.Reader, offset, total int64) error {
	buf := make([]byte, 256*1024)
	written := int64(0)
	start := time.Now()
	var lastPrint time.Time

	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, werr := dst.Write(buf[:n]); werr != nil {
				return fmt.Errorf("write: %w", werr)
			}

			written += int64(n)

			if now := time.Now(); now.Sub(lastPrint) >= 250*time.Millisecond {
				lastPrint = now
				printProgress(offset+written, written, total, time.Since(start), false)
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("read: %w", err)
		}
	}

	printProgress(offset+written, written, total, time.Since(start), true)

	return nil
}

func printProgress(done, written, total int64, elapsed time.Duration, final bool) {
	secs := max(elapsed.Seconds(), 0.001)
	speed := float64(written) / secs / 1024 / 1024 //nolint:mnd

	if total > 0 {
		fmt.Printf("\r  %.1f%% (%s / %s) @ %.1f MB/s   ",
			float64(done)/float64(total)*100, fmtBytes(done), fmtBytes(total), speed)
	} else {
		fmt.Printf("\r  %s @ %.1f MB/s   ", fmtBytes(done), speed)
	}

	if final {
		fmt.Println()
	}
}

func fmtBytes(n int64) string {
	const (
		kb = 1 << 10
		mb = 1 << 20
		gb = 1 << 30
	)

	switch {
	case n >= gb:
		return fmt.Sprintf("%.2f GB", float64(n)/gb)
	case n >= mb:
		return fmt.Sprintf("%.1f MB", float64(n)/mb)
	case n >= kb:
		return fmt.Sprintf("%.0f KB", float64(n)/kb)
	default:
		return fmt.Sprintf("%d B", n)
	}
}
