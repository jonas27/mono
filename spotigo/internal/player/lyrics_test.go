package player

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func newLyrics(syncType string, lines ...[2]string) *lyricsResponse {
	lr := &lyricsResponse{}
	lr.Lyrics.SyncType = syncType
	for _, l := range lines {
		lr.Lyrics.Lines = append(lr.Lyrics.Lines, struct {
			StartTimeMs string `json:"startTimeMs"`
			Words       string `json:"words"`
		}{StartTimeMs: l[0], Words: l[1]})
	}
	return lr
}

func TestFormatLRC(t *testing.T) {
	lr := newLyrics("LINE_SYNCED",
		[2]string{"0", "First line"},
		[2]string{"61230", "Past a minute"},
		[2]string{"bad", "Unparseable time"},
	)

	got := formatLRC(lr)
	want := "[00:00.00]First line\n[01:01.23]Past a minute\n[00:00.00]Unparseable time\n"
	require.Equal(t, want, got)
}

func TestFormatPlain(t *testing.T) {
	lr := newLyrics("UNSYNCED",
		[2]string{"0", "Line one"},
		[2]string{"0", "Line two"},
	)

	require.Equal(t, "Line one\nLine two\n", formatPlain(lr))
}
