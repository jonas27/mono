// Package config handles loading and saving application configuration.
package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const configFile = "config.json"

// APIKeys holds the API credentials for external services.
type APIKeys struct {
	SpotifyID     string `json:"spotify_id"`
	SpotifySecret string `json:"spotify_secret"`
	GeniusToken   string `json:"genius_token"`
}

// Config represents the application configuration stored in config.json.
type Config struct {
	APIKeys      APIKeys `json:"api_keys"`
	DownloadPath string  `json:"download_path"`
	Quality      string  `json:"quality"`
	SourceType   string  `json:"source_type"`
}

// Load reads the configuration from config.json in the working directory.
// If the file does not exist or cannot be parsed, a default config is returned.
func Load() Config {
	cfg := Config{
		DownloadPath: "downloads",
		Quality:      "192",
		SourceType:   "Spotify",
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return cfg
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg
	}

	if cfg.DownloadPath == "" {
		cfg.DownloadPath = "downloads"
	}

	if cfg.Quality == "" {
		cfg.Quality = "192"
	}

	return cfg
}

// Save writes the configuration to config.json in the working directory.
func Save(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0o600); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

// PromptForKeys interactively prompts the user to enter API keys via stdin.
func PromptForKeys(current APIKeys) (APIKeys, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("--- API Key Setup ---")
	fmt.Println("Genius.com token is optional (for lyrics).")

	fmt.Printf("Spotify Client ID [%s]: ", mask(current.SpotifyID))
	spotifyID, err := reader.ReadString('\n')
	if err != nil {
		return current, fmt.Errorf("read spotify id: %w", err)
	}

	spotifyID = strings.TrimSpace(spotifyID)
	if spotifyID == "" {
		spotifyID = current.SpotifyID
	}

	fmt.Printf("Spotify Client Secret [%s]: ", mask(current.SpotifySecret))
	spotifySecret, err := reader.ReadString('\n')
	if err != nil {
		return current, fmt.Errorf("read spotify secret: %w", err)
	}

	spotifySecret = strings.TrimSpace(spotifySecret)
	if spotifySecret == "" {
		spotifySecret = current.SpotifySecret
	}

	fmt.Printf("Genius.com Access Token (optional) [%s]: ", mask(current.GeniusToken))
	geniusToken, err := reader.ReadString('\n')
	if err != nil {
		return current, fmt.Errorf("read genius token: %w", err)
	}

	geniusToken = strings.TrimSpace(geniusToken)
	if geniusToken == "" {
		geniusToken = current.GeniusToken
	}

	return APIKeys{
		SpotifyID:     spotifyID,
		SpotifySecret: spotifySecret,
		GeniusToken:   geniusToken,
	}, nil
}

// mask returns the first 4 chars of s followed by "****", or empty if s is empty.
func mask(s string) string {
	if s == "" {
		return ""
	}

	if len(s) <= 4 {
		return "****"
	}

	return s[:4] + "****"
}
