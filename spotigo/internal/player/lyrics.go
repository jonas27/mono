package player

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	librespot "github.com/devgianlu/go-librespot"
)

// lyricsResponse mirrors the JSON returned by Spotify's color-lyrics endpoint.
type lyricsResponse struct {
	Lyrics struct {
		SyncType string `json:"syncType"` // "LINE_SYNCED" or "UNSYNCED"
		Lines    []struct {
			StartTimeMs string `json:"startTimeMs"`
			Words       string `json:"words"`
		} `json:"lines"`
		Language string `json:"language"`
	} `json:"lyrics"`
}

// downloadLyrics fetches lyrics for id and writes them next to the track in dir
// as <title>.lrc (line-synced with timestamps, or plain lines when unsynced).
// Errors and missing lyrics are logged and swallowed — lyrics are best-effort.
func (s *Session) downloadLyrics(ctx context.Context, id librespot.SpotifyId, dir string, meta trackMeta) {
	path := filepath.Join(dir, safeFilename(meta.Title)) + ".lrc"
	if _, err := os.Stat(path); err == nil {
		return // already have lyrics for this track
	}

	lr, err := s.fetchLyrics(ctx, id)
	if err != nil {
		slog.Warn("fetch lyrics failed", "track", meta.Title, "error", err)
		return
	}
	if lr == nil || len(lr.Lyrics.Lines) == 0 {
		return // no lyrics available for this track
	}

	// Synced lyrics get [mm:ss.cc] timestamps; unsynced ones are written as
	// plain lines in the same .lrc file (players accept both).
	content := formatPlain(lr)
	if lr.Lyrics.SyncType == "LINE_SYNCED" {
		content = formatLRC(lr)
	}
	if err := os.WriteFile(path, []byte(content), 0o640); err != nil {
		slog.Warn("write lyrics failed", "file", filepath.Base(path), "error", err)
		return
	}
	slog.Info("wrote lyrics", "file", filepath.Base(path))
}

// fetchLyrics requests lyrics for id from the spclient color-lyrics endpoint.
// Returns nil (no error) when the track has no lyrics.
func (s *Session) fetchLyrics(ctx context.Context, id librespot.SpotifyId) (*lyricsResponse, error) {
	resp, err := s.sess.Spclient().Request(ctx, http.MethodGet,
		"/color-lyrics/v2/track/"+id.Base62(),
		url.Values{"format": {"json"}, "vocalRemoval": {"false"}, "market": {"from_token"}},
		// App-Platform is required or the endpoint returns 400.
		http.Header{"App-Platform": {"WebPlayer"}},
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // no lyrics for this track
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lyrics request returned status %d", resp.StatusCode)
	}

	var lr lyricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		return nil, fmt.Errorf("decode lyrics: %w", err)
	}
	return &lr, nil
}

// formatLRC renders line-synced lyrics as an LRC file with [mm:ss.cc] timestamps.
func formatLRC(lr *lyricsResponse) string {
	var b strings.Builder
	for _, ln := range lr.Lyrics.Lines {
		ms, err := strconv.Atoi(ln.StartTimeMs)
		if err != nil {
			ms = 0
		}
		m := ms / 60000
		sec := ms / 1000 % 60
		cs := ms / 10 % 100
		fmt.Fprintf(&b, "[%02d:%02d.%02d]%s\n", m, sec, cs, ln.Words)
	}
	return b.String()
}

// formatPlain renders unsynced lyrics as plain text, one line per lyric line.
func formatPlain(lr *lyricsResponse) string {
	var b strings.Builder
	for _, ln := range lr.Lyrics.Lines {
		b.WriteString(ln.Words)
		b.WriteByte('\n')
	}
	return b.String()
}
