package player

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeOpusFromFile(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not in PATH")
	}

	dir := t.TempDir()

	// 0.1 s of silence: 44100 Hz, 2 ch, f32le → 4 bytes/sample
	silence := make([]byte, 44100/10*2*4)
	rawPath := filepath.Join(dir, "silence.raw")
	require.NoError(t, os.WriteFile(rawPath, silence, 0o600))

	outPath := filepath.Join(dir, "out.opus")
	meta := trackMeta{Title: "Test Track", Artist: "Test Artist", Album: "Test Album", TrackNumber: 1}

	require.NoError(t, encodeOpusFromFile(context.Background(), rawPath, outPath, meta))

	info, err := os.Stat(outPath)
	require.NoError(t, err, "output file not created")
	require.Positive(t, info.Size(), "output file is empty")

	out, err := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json",
		"-show_streams", outPath).Output()
	require.NoError(t, err)

	var probe struct {
		Streams []struct {
			Tags map[string]string `json:"tags"`
		} `json:"streams"`
	}
	require.NoError(t, json.Unmarshal(out, &probe))
	require.NotEmpty(t, probe.Streams)

	require.Equal(t, "1", probe.Streams[0].Tags["track"])
}

func TestEncodeOpusDiscNumber(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not in PATH")
	}

	dir := t.TempDir()
	silence := make([]byte, 44100/10*2*4)
	rawPath := filepath.Join(dir, "silence.raw")
	require.NoError(t, os.WriteFile(rawPath, silence, 0o600))

	outPath := filepath.Join(dir, "out.opus")
	meta := trackMeta{Title: "Dani California", Artist: "Red Hot Chili Peppers", Album: "Stadium Arcadium", TrackNumber: 1, DiscNumber: 1}
	require.NoError(t, encodeOpusFromFile(context.Background(), rawPath, outPath, meta))

	out, err := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json",
		"-show_streams", outPath).Output()
	require.NoError(t, err)

	var probe struct {
		Streams []struct {
			Tags map[string]string `json:"tags"`
		} `json:"streams"`
	}
	require.NoError(t, json.Unmarshal(out, &probe))
	require.NotEmpty(t, probe.Streams)

	tags := probe.Streams[0].Tags
	require.Equal(t, "1", tags["disc"])
	require.Equal(t, "1", tags["track"])
}

func TestEncodeOpusFromFileMissingInput(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not in PATH")
	}

	dir := t.TempDir()
	err := encodeOpusFromFile(context.Background(), filepath.Join(dir, "nonexistent.raw"), filepath.Join(dir, "out.opus"), trackMeta{})
	require.Error(t, err)
}
