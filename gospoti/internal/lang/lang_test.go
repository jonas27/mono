//go:build unit

package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet_KnownKeyAllLanguages(t *testing.T) {
	tests := []struct {
		key      string
		lang     Code
		expected string
	}{
		{"window_title", EN, "SpotiDown"},
		{"window_title", BG, "SpotiDown"},
		{"window_title", ES, "SpotiDown"},
		{"start_button", EN, "Start Download"},
		{"start_button", BG, "Старт"},
		{"start_button", ES, "Iniciar Descarga"},
		{"cancel_button", EN, "Cancel"},
		{"cancel_button", BG, "Отказ"},
		{"cancel_button", ES, "Cancelar"},
		{"downloading_song", EN, "Downloading:"},
		{"downloading_song", BG, "Сваля се:"},
		{"downloading_song", ES, "Descargando:"},
		{"song_done", EN, "Done:"},
		{"song_done", BG, "Готово:"},
		{"song_done", ES, "Hecho:"},
	}

	for _, tt := range tests {
		t.Run(tt.key+"/"+string(tt.lang), func(t *testing.T) {
			got := Get(tt.key, tt.lang)
			require.NotEmpty(t, got)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGet_UnknownLangFallsBackToEnglish(t *testing.T) {
	got := Get("start_button", Code("xx"))
	assert.Equal(t, "Start Download", got, "unknown lang code should fall back to English")
}

func TestGet_UnknownKeyReturnsBracketedKey(t *testing.T) {
	got := Get("no_such_key", EN)
	assert.Equal(t, "[no_such_key]", got)
}

func TestGet_UnknownKeyUnknownLangReturnsBracketedKey(t *testing.T) {
	got := Get("missing", Code("fr"))
	assert.Equal(t, "[missing]", got)
}

func TestGet_AllKeysHaveEnglishValue(t *testing.T) {
	// Every key registered in the strings map must resolve to a non-empty English string.
	for key := range strings {
		got := Get(key, EN)
		assert.NotEqual(t, "["+key+"]", got, "key %q should have an English value", key)
		assert.NotEmpty(t, got, "key %q English value must not be empty", key)
	}
}

func TestGet_AllKeysAllSupportedLanguages(t *testing.T) {
	langs := []Code{EN, BG, ES}
	for key := range strings {
		for _, l := range langs {
			got := Get(key, l)
			assert.NotEmpty(t, got, "key %q lang %q must return a non-empty string", key, l)
		}
	}
}
