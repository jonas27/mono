// Package scraper fetches and parses tiz-cycling.tv pages.
package scraper

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
)

var (
	// reTape1 extracts the direct MP4 URL from the video.php iframe src on tape-1 pages.
	reTape1 = regexp.MustCompile(`video\.php\?v=([^"]+)"`)

	// reCatVideo extracts video page URLs from the <h3> links in a category page's main content.
	reCatVideo = regexp.MustCompile(`<h3><a\s+href="(https://tiz-cycling\.tv/video/[^"]+)"`)

	// reTitle extracts the page title from the <title> element.
	reTitle = regexp.MustCompile(`<title>([^<]+)</title>`)
)

// VideoInfo holds the extracted details needed to download a single video.
type VideoInfo struct {
	Title     string
	DirectURL string
	Filename  string
}

// Scraper fetches and parses tiz-cycling.tv pages.
type Scraper struct {
	client *http.Client
}

// New returns a Scraper ready to use.
func New() *Scraper {
	return &Scraper{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Resolve returns video page URLs for the given URL.
// Category pages (/categories/...) are expanded into their member video URLs.
// Video pages (/video/...) are returned as-is.
func (s *Scraper) Resolve(ctx context.Context, rawURL string) ([]string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	switch {
	case strings.HasPrefix(u.Path, "/categories/"):
		return s.fetchCategory(ctx, rawURL)
	case strings.HasPrefix(u.Path, "/video/"):
		return []string{rawURL}, nil
	default:
		return nil, fmt.Errorf("expected /video/... or /categories/... path, got: %s", u.Path)
	}
}

// FetchVideo returns the VideoInfo for a single video page URL.
// It requires the page to have a tape-1 direct MP4 player; other mirrors are not supported.
func (s *Scraper) FetchVideo(ctx context.Context, pageURL string) (VideoInfo, error) {
	body, err := s.fetchBody(ctx, pageURL)
	if err != nil {
		return VideoInfo{}, err
	}

	m := reTape1.FindSubmatch(body)
	if m == nil {
		return VideoInfo{}, fmt.Errorf("no direct MP4 mirror found on %s", pageURL)
	}

	directURL := html.UnescapeString(string(m[1]))
	title := extractTitle(body)

	filename, err := urlToFilename(directURL)
	if err != nil {
		filename = sanitize(title) + ".mp4"
	}

	return VideoInfo{
		Title:     title,
		DirectURL: directURL,
		Filename:  filename,
	}, nil
}

func (s *Scraper) fetchCategory(ctx context.Context, categoryURL string) ([]string, error) {
	body, err := s.fetchBody(ctx, categoryURL)
	if err != nil {
		return nil, err
	}

	matches := reCatVideo.FindAllSubmatch(body, -1)
	seen := make(map[string]bool, len(matches))
	urls := make([]string, 0, len(matches))

	for _, m := range matches {
		u := string(m[1])
		if !seen[u] {
			seen[u] = true
			urls = append(urls, u)
		}
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("no videos found in category: %s", categoryURL)
	}

	return urls, nil
}

func (s *Scraper) fetchBody(ctx context.Context, u string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; gotiz/1.0)")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: HTTP %d", u, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func extractTitle(body []byte) string {
	m := reTitle.FindSubmatch(body)
	if m == nil {
		return "unknown"
	}

	t := html.UnescapeString(strings.TrimSpace(string(m[1])))
	t = strings.TrimSuffix(t, " - Tiz-Cycling")
	t = strings.TrimSuffix(t, " – Tiz-Cycling") // en dash variant

	return t
}

// urlToFilename derives a human-readable filename from a video.tiz-cycling.io URL.
// Backblaze B2 encodes spaces as + in the path, so we reverse that.
func urlToFilename(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse URL: %w", err)
	}

	name := path.Base(u.Path)
	if name == "" || name == "." || name == "/" {
		return "", fmt.Errorf("cannot derive filename from path: %s", u.Path)
	}

	return strings.ReplaceAll(name, "+", " "), nil
}

func sanitize(s string) string {
	var b strings.Builder

	for _, r := range s {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			b.WriteRune('_')
		default:
			b.WriteRune(r)
		}
	}

	return b.String()
}
