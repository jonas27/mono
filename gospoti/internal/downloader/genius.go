// Package downloader handles music downloading, metadata enrichment and yt-dlp integration.
package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var lyricsTagRe = regexp.MustCompile(`<[^>]+>`)

// GeniusClient fetches song lyrics from the Genius API.
type GeniusClient struct {
	token  string
	client *http.Client
}

// NewGeniusClient creates a new GeniusClient using the provided API token.
func NewGeniusClient(token string) *GeniusClient {
	return &GeniusClient{
		token:  token,
		client: &http.Client{},
	}
}

type geniusSearchResponse struct {
	Response struct {
		Hits []struct {
			Result struct {
				URL string `json:"url"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

// GetLyrics searches Genius for the given artist+title and returns plain-text lyrics.
// Returns an empty string if lyrics are not found.
func (g *GeniusClient) GetLyrics(ctx context.Context, artist, title string) (string, error) {
	query := url.QueryEscape(artist + " " + title)
	searchURL := "https://api.genius.com/search?q=" + query

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return "", fmt.Errorf("create Genius search request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.token)

	resp, err := g.client.Do(req) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("genius search request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read Genius search response: %w", err)
	}

	var result geniusSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parse Genius search response: %w", err)
	}

	if len(result.Response.Hits) == 0 {
		return "", nil
	}

	lyricsPageURL := result.Response.Hits[0].Result.URL
	if lyricsPageURL == "" {
		return "", nil
	}

	return g.scrapeLyricsPage(ctx, lyricsPageURL)
}

func (g *GeniusClient) scrapeLyricsPage(ctx context.Context, pageURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return "", fmt.Errorf("create Genius page request: %w", err)
	}

	resp, err := g.client.Do(req) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("fetch Genius lyrics page: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read Genius lyrics page: %w", err)
	}

	html := string(body)

	// Find all <div data-lyrics-container="true"> blocks.
	var lyricsBlocks []string
	remaining := html

	for {
		start := strings.Index(remaining, `data-lyrics-container="true"`)
		if start == -1 {
			break
		}

		// Find the opening tag start.
		tagStart := strings.LastIndex(remaining[:start], "<")
		if tagStart == -1 {
			break
		}

		// Find the matching closing </div>.
		depth := 0
		pos := tagStart

		for pos < len(remaining) {
			if remaining[pos] == '<' {
				if pos+1 < len(remaining) && remaining[pos+1] == '/' {
					depth--
					if depth == 0 {
						closeEnd := strings.Index(remaining[pos:], ">")
						if closeEnd != -1 {
							lyricsBlocks = append(lyricsBlocks, remaining[tagStart:pos+closeEnd+1])
						}

						remaining = remaining[pos+1:]
						break
					}
				} else if remaining[pos] != '<' || remaining[pos+1] != '!' {
					depth++
				}
			}

			pos++
		}

		if pos >= len(remaining) {
			break
		}
	}

	if len(lyricsBlocks) == 0 {
		return "", nil
	}

	// Strip HTML tags and clean up the lyrics text.
	var parts []string

	for _, block := range lyricsBlocks {
		// Replace <br> tags with newlines before stripping.
		block = regexp.MustCompile(`(?i)<br\s*/?>`).ReplaceAllString(block, "\n")
		// Strip remaining HTML tags.
		clean := lyricsTagRe.ReplaceAllString(block, "")
		parts = append(parts, clean)
	}

	lyrics := strings.Join(parts, "\n")
	lyrics = strings.TrimSpace(lyrics)

	return lyrics, nil
}
