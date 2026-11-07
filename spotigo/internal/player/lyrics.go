package player

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// lrclibBaseURL is the auth-free LRCLIB lyrics API (https://lrclib.net).
const lrclibBaseURL = "https://lrclib.net"

// lrclibUserAgent identifies spotigo to LRCLIB, as their API guidelines request.
const lrclibUserAgent = "spotigo (https://github.com/jmanser/spotigo)"

// lrclibRecord is a single lyrics record returned by the LRCLIB API.
type lrclibRecord struct {
	Instrumental bool   `json:"instrumental"`
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"syncedLyrics"` // LRC format with [mm:ss.xx] timestamps
}

// downloadLyrics fetches lyrics for meta from LRCLIB and writes them next to the
// track in dir as <title>.lrc (synced when available, otherwise plain lines).
// Errors and missing lyrics are logged and swallowed — lyrics are best-effort.
func (s *Session) downloadLyrics(ctx context.Context, dir string, meta trackMeta) {
	if meta.Title == "" || meta.Artist == "" {
		return
	}
	path := filepath.Join(dir, safeFilename(meta.Title)) + ".lrc"
	if _, err := os.Stat(path); err == nil {
		return // already have lyrics for this track
	}

	rec, err := fetchLyrics(ctx, meta)
	if err != nil {
		slog.Warn("fetch lyrics failed", "track", meta.Title, "error", err)
		return
	}
	if rec == nil || rec.Instrumental {
		return // no lyrics available for this track
	}

	content := rec.SyncedLyrics
	if strings.TrimSpace(content) == "" {
		content = rec.PlainLyrics
	}
	if strings.TrimSpace(content) == "" {
		return
	}
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	if err := os.WriteFile(path, []byte(content), 0o640); err != nil {
		slog.Warn("write lyrics failed", "file", filepath.Base(path), "error", err)
		return
	}
	slog.Info("wrote lyrics", "file", filepath.Base(path))
}

// fetchLyrics looks up lyrics for meta on LRCLIB. It tries an exact match first
// (artist + track + album + duration), falling back to a fuzzy search. Returns
// nil (no error) when no lyrics are found.
func fetchLyrics(ctx context.Context, meta trackMeta) (*lrclibRecord, error) {
	if meta.Duration > 0 {
		rec, err := lrclibGet(ctx, meta)
		if err != nil {
			return nil, err
		}
		if rec != nil {
			return rec, nil
		}
		// exact match missed — fall through to fuzzy search
	}
	return lrclibSearch(ctx, meta)
}

// lrclibGet performs an exact LRCLIB lookup. Returns nil when there is no match.
func lrclibGet(ctx context.Context, meta trackMeta) (*lrclibRecord, error) {
	q := url.Values{
		"artist_name": {meta.Artist},
		"track_name":  {meta.Title},
		"duration":    {strconv.Itoa(meta.Duration)},
	}
	if meta.Album != "" {
		q.Set("album_name", meta.Album)
	}

	body, status, err := lrclibRequest(ctx, "/api/get", q)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("lrclib get returned status %d", status)
	}
	var rec lrclibRecord
	if err := json.Unmarshal(body, &rec); err != nil {
		return nil, fmt.Errorf("decode lrclib get: %w", err)
	}
	return &rec, nil
}

// lrclibSearch performs a fuzzy LRCLIB search and returns the best match,
// preferring a result that has synced lyrics. Returns nil when nothing matches.
func lrclibSearch(ctx context.Context, meta trackMeta) (*lrclibRecord, error) {
	q := url.Values{
		"artist_name": {meta.Artist},
		"track_name":  {meta.Title},
	}

	body, status, err := lrclibRequest(ctx, "/api/search", q)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("lrclib search returned status %d", status)
	}
	var results []lrclibRecord
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("decode lrclib search: %w", err)
	}
	return pickLyrics(results), nil
}

// pickLyrics chooses the best usable record from search results: the first with
// synced lyrics, else the first with plain lyrics. Instrumental records and
// records with no lyrics are ignored. Returns nil if none are usable.
func pickLyrics(results []lrclibRecord) *lrclibRecord {
	var fallback *lrclibRecord
	for i := range results {
		r := &results[i]
		if r.Instrumental {
			continue
		}
		if strings.TrimSpace(r.SyncedLyrics) != "" {
			return r // prefer synced lyrics
		}
		if fallback == nil && strings.TrimSpace(r.PlainLyrics) != "" {
			fallback = r
		}
	}
	return fallback
}

// lrclibRequest issues a GET to the LRCLIB API and returns the body and status.
func lrclibRequest(ctx context.Context, path string, q url.Values) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, lrclibBaseURL+path+"?"+q.Encode(), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", lrclibUserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}
