//go:build unit

package downloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// cleanFilename
// ---------------------------------------------------------------------------

func TestCleanFilename_RemovesUnsafeChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Safe characters are preserved as-is.
		{"Artist - Title", "Artist - Title"},
		// Colon is stripped; adjacent spaces collapse only if trimmed at edges.
		{"Artist: Title", "Artist Title"},
		// Slash is stripped; no space inserted.
		{"Hello/World", "HelloWorld"},
		{"AC/DC - Back in Black", "ACDC - Back in Black"},
		// Parentheses and dot stripped.
		{"Song (feat. Artist)", "Song feat Artist"},
		// Hash stripped.
		{"Track #1", "Track 1"},
		// Leading/trailing spaces trimmed.
		{"  leading trailing  ", "leading trailing"},
		// Unicode letters (CJK, umlauts) are kept by unicode.IsLetter.
		{"日本語 Title", "日本語 Title"},
		{"Ärger - Über", "Ärger - Über"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := cleanFilename(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestCleanFilename_KeepsLettersDigitsSpacesDashes(t *testing.T) {
	input := "Artist123 - Track456"
	got := cleanFilename(input)
	assert.Equal(t, input, got)
}

func TestCleanFilename_EmptyInput(t *testing.T) {
	assert.Equal(t, "", cleanFilename(""))
}

// ---------------------------------------------------------------------------
// spotifyIDFromURL
// ---------------------------------------------------------------------------

func TestSpotifyIDFromURL(t *testing.T) {
	tests := []struct {
		rawURL   string
		expected string
	}{
		{"https://open.spotify.com/track/4iV5W9uYEdYUVa79Axb7Rh", "4iV5W9uYEdYUVa79Axb7Rh"},
		{"https://open.spotify.com/playlist/37i9dQZF1DXcBWIGoYBM5M?si=abc", "37i9dQZF1DXcBWIGoYBM5M"},
		{"https://open.spotify.com/album/6dVIqQ8qmQ5GBnJ9shOYGE", "6dVIqQ8qmQ5GBnJ9shOYGE"},
		{"https://example.com/something-else", ""},
	}

	for _, tt := range tests {
		t.Run(tt.rawURL, func(t *testing.T) {
			got := spotifyIDFromURL(tt.rawURL)
			assert.Equal(t, tt.expected, got)
		})
	}
}
