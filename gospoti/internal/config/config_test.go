//go:build unit

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// withTempDir changes the working directory to a fresh temp dir for the duration
// of the test so that Load/Save operate on isolated files.
func withTempDir(t *testing.T) {
	t.Helper()
	orig, err := os.Getwd()
	require.NoError(t, err)
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(orig) })
}

// TestLoad_DefaultsWhenFileMissing verifies Load returns sensible defaults when
// config.json does not exist.
func TestLoad_DefaultsWhenFileMissing(t *testing.T) {
	withTempDir(t)

	cfg := Load()

	assert.Equal(t, "downloads", cfg.DownloadPath)
	assert.Equal(t, "192", cfg.Quality)
	assert.Equal(t, "Spotify", cfg.SourceType)
	assert.Empty(t, cfg.APIKeys.SpotifyID)
	assert.Empty(t, cfg.APIKeys.SpotifySecret)
	assert.Empty(t, cfg.APIKeys.GeniusToken)
}

// TestLoad_ParsesValidJSON confirms all fields are read correctly.
func TestLoad_ParsesValidJSON(t *testing.T) {
	withTempDir(t)

	input := Config{
		APIKeys: APIKeys{
			SpotifyID:     "my-id",
			SpotifySecret: "my-secret",
			GeniusToken:   "genius-tok",
		},
		DownloadPath: "/music",
		Quality:      "320",
		SourceType:   "Apple Music",
	}

	data, err := json.MarshalIndent(input, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile("config.json", data, 0o600))

	cfg := Load()

	assert.Equal(t, "my-id", cfg.APIKeys.SpotifyID)
	assert.Equal(t, "my-secret", cfg.APIKeys.SpotifySecret)
	assert.Equal(t, "genius-tok", cfg.APIKeys.GeniusToken)
	assert.Equal(t, "/music", cfg.DownloadPath)
	assert.Equal(t, "320", cfg.Quality)
	assert.Equal(t, "Apple Music", cfg.SourceType)
}

// TestLoad_MalformedJSONReturnsDefault confirms Load falls back to defaults on
// parse failure instead of returning an error.
func TestLoad_MalformedJSONReturnsDefault(t *testing.T) {
	withTempDir(t)

	require.NoError(t, os.WriteFile("config.json", []byte("{not valid json"), 0o600))

	cfg := Load()

	assert.Equal(t, "downloads", cfg.DownloadPath)
	assert.Equal(t, "192", cfg.Quality)
}

// TestLoad_EmptyDownloadPathFallsBackToDefault checks the post-unmarshal
// fallback for empty DownloadPath.
func TestLoad_EmptyDownloadPathFallsBackToDefault(t *testing.T) {
	withTempDir(t)

	raw := `{"download_path":"","quality":"256","source_type":"Spotify"}`
	require.NoError(t, os.WriteFile("config.json", []byte(raw), 0o600))

	cfg := Load()

	assert.Equal(t, "downloads", cfg.DownloadPath)
	assert.Equal(t, "256", cfg.Quality)
}

// TestLoad_EmptyQualityFallsBackToDefault checks the post-unmarshal fallback
// for empty Quality.
func TestLoad_EmptyQualityFallsBackToDefault(t *testing.T) {
	withTempDir(t)

	raw := `{"download_path":"/music","quality":"","source_type":"Spotify"}`
	require.NoError(t, os.WriteFile("config.json", []byte(raw), 0o600))

	cfg := Load()

	assert.Equal(t, "192", cfg.Quality)
	assert.Equal(t, "/music", cfg.DownloadPath)
}

// TestSave_RoundTrip verifies that saving and loading a config produces the
// original values.
func TestSave_RoundTrip(t *testing.T) {
	withTempDir(t)

	original := Config{
		APIKeys:      APIKeys{SpotifyID: "id", SpotifySecret: "secret"},
		DownloadPath: "/tmp/music",
		Quality:      "128",
		SourceType:   "YouTube",
	}

	require.NoError(t, Save(original))

	loaded := Load()
	assert.Equal(t, original, loaded)
}

// TestSave_CreatesConfigFile checks that the file is actually written.
func TestSave_CreatesConfigFile(t *testing.T) {
	withTempDir(t)

	require.NoError(t, Save(Config{DownloadPath: "out", Quality: "320"}))

	_, err := os.Stat(filepath.Join("config.json"))
	assert.NoError(t, err)
}

// TestMask_* tests the unexported mask helper via table-driven cases.
func TestMask(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"abc", "****"},
		{"abcd", "****"},
		{"abcde", "abcd****"},
		{"verylongtoken", "very****"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, mask(tt.input))
		})
	}
}
