package player

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPickLyricsPrefersSynced(t *testing.T) {
	results := []lrclibRecord{
		{PlainLyrics: "plain only"},
		{SyncedLyrics: "[00:01.00]synced", PlainLyrics: "also plain"},
	}
	got := pickLyrics(results)
	require.NotNil(t, got)
	require.Equal(t, "[00:01.00]synced", got.SyncedLyrics)
}

func TestPickLyricsFallsBackToPlain(t *testing.T) {
	results := []lrclibRecord{
		{Instrumental: true, PlainLyrics: "ignored"},
		{PlainLyrics: "first plain"},
		{PlainLyrics: "second plain"},
	}
	got := pickLyrics(results)
	require.NotNil(t, got)
	require.Equal(t, "first plain", got.PlainLyrics)
}

func TestPickLyricsSkipsInstrumentalAndEmpty(t *testing.T) {
	results := []lrclibRecord{
		{Instrumental: true, SyncedLyrics: "[00:01.00]nope"},
		{},                   // no lyrics
		{PlainLyrics: "   "}, // whitespace only
	}
	require.Nil(t, pickLyrics(results))
}

func TestPickLyricsEmpty(t *testing.T) {
	require.Nil(t, pickLyrics(nil))
}
