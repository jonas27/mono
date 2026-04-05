// Package lang provides localised string lookup for en, bg and es.
package lang

// Code represents a supported language code.
type Code string

const (
	// EN is the English language code.
	EN Code = "en"
	// BG is the Bulgarian language code.
	BG Code = "bg"
	// ES is the Spanish language code.
	ES Code = "es"
)

var strings = map[string]map[Code]string{
	"window_title": {EN: "SpotiDown", BG: "SpotiDown", ES: "SpotiDown"},
	"url_label": {
		EN: "URL (Spotify or Apple Music):",
		BG: "URL (Spotify или Apple Music):",
		ES: "URL (Spotify o Apple Music):",
	},
	"start_button":   {EN: "Start Download", BG: "Старт", ES: "Iniciar Descarga"},
	"progress_label": {EN: "Log:", BG: "Прогрес:", ES: "Registro:"},
	"lang_label":     {EN: "Language:", BG: "Език:", ES: "Idioma:"},
	"source_type_label": {
		EN: "Source Type:",
		BG: "Тип на източника:",
		ES: "Tipo de Fuente:",
	},
	"download_folder_label": {
		EN: "Download Folder:",
		BG: "Папка за сваляне:",
		ES: "Carpeta de Descarga:",
	},
	"browse_button":     {EN: "Browse...", BG: "Избери...", ES: "Explorar..."},
	"keys_prompt_title": {EN: "API Keys Required", BG: "Необходими са API Ключове", ES: "Se Requieren Claves de API"},
	"keys_prompt_message": {
		EN: "Please enter your API keys. Genius.com key is optional.",
		BG: "Моля, въведете вашите API ключове. Ключът за Genius.com не е задължителен.",
		ES: "Por favor, ingrese sus claves de API. La clave de Genius.com es opcional.",
	},
	"save_keys_button": {EN: "Save and Continue", BG: "Запази и Продължи", ES: "Guardar y Continuar"},
	"settings_menu_label": {
		EN: "Settings",
		BG: "Настройки",
		ES: "Configuración",
	},
	"settings_menu_change_keys": {
		EN: "Change API Keys",
		BG: "Смяна на API Ключове",
		ES: "Cambiar Claves de API",
	},
	"quality_label": {
		EN: "Audio Quality (kbps):",
		BG: "Качество на звука (kbps):",
		ES: "Calidad de Audio (kbps):",
	},
	"cancel_button": {EN: "Cancel", BG: "Отказ", ES: "Cancelar"},
	"cancelling_message": {
		EN: "Cancelling... will stop after the current song.",
		BG: "Отменя се... ще спре след текущата песен.",
		ES: "Cancelando... se detendrá después de la canción actual.",
	},
	"downloading_song": {EN: "Downloading:", BG: "Сваля се:", ES: "Descargando:"},
	"song_done":        {EN: "Done:", BG: "Готово:", ES: "Hecho:"},
	"song_skipped": {
		EN: "Skipped (file already exists):",
		BG: "Пропуснато (файлът вече съществува):",
		ES: "Omitido (el archivo ya existe):",
	},
	"process_finished": {
		EN: "\n--- Finished! ---",
		BG: "\n--- Процесът приключи! ---",
		ES: "\n--- ¡Proceso Terminado! ---",
	},
	"connecting_spotify": {
		EN: "Connecting to Spotify API...",
		BG: "Свързване със Spotify API...",
		ES: "Conectando a la API de Spotify...",
	},
	"connecting_genius": {
		EN: "Connecting to Genius.com API...",
		BG: "Свързване с Genius.com API...",
		ES: "Conectando a la API de Genius.com...",
	},
	"found_tracks": {
		EN: "Found %d tracks.",
		BG: "Намерени %d песни.",
		ES: "Se encontraron %d pistas.",
	},
}

// Get returns the localised string for the given key and language code.
// Falls back to English, then to the key itself if not found.
func Get(key string, lang Code) string {
	langMap, ok := strings[key]
	if !ok {
		return "[" + key + "]"
	}

	if val, ok := langMap[lang]; ok {
		return val
	}

	if val, ok := langMap[EN]; ok {
		return val
	}

	return "[" + key + "]"
}
