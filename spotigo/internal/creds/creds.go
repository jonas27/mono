// Package creds handles loading and saving Spotify credentials from a YAML file.
package creds

import (
	"encoding/base64"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Creds holds Spotify authentication credentials persisted to .creds.yaml.
type Creds struct {
	Username string `yaml:"username,omitempty"`
	Stored   string `yaml:"stored_credentials,omitempty"` // base64-encoded blob
}

// Load reads credentials from path. Returns empty Creds if the file doesn't exist.
func Load(path string) (*Creds, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Creds{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var c Creds
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &c, nil
}

// Save writes credentials to path with mode 0600.
func Save(path string, c *Creds) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// StoredBytes decodes the base64 blob, returning nil if empty.
func (c *Creds) StoredBytes() ([]byte, error) {
	if c.Stored == "" {
		return nil, nil
	}
	b, err := base64.StdEncoding.DecodeString(c.Stored)
	if err != nil {
		return nil, fmt.Errorf("decode stored_credentials: %w", err)
	}
	return b, nil
}

// SetStored encodes blob to base64 and stores it.
func (c *Creds) SetStored(blob []byte) {
	c.Stored = base64.StdEncoding.EncodeToString(blob)
}
