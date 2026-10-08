package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

func path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("find config dir: %w", err)
	}
	return filepath.Join(dir, "worklog", "config.json"), nil
}

// Load reads the saved config; WL_URL and WL_TOKEN override it.
func Load() (Config, error) {
	var cfg Config

	p, err := path()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(p)
	switch {
	case err == nil:
		if err := json.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse %s: %w", p, err)
		}
	case !errors.Is(err, os.ErrNotExist):
		return Config{}, fmt.Errorf("read %s: %w", p, err)
	}

	if v := os.Getenv("WL_URL"); v != "" {
		cfg.URL = v
	}
	if v := os.Getenv("WL_TOKEN"); v != "" {
		cfg.Token = v
	}

	cfg.URL = strings.TrimRight(cfg.URL, "/")
	if cfg.URL == "" || cfg.Token == "" {
		return Config{}, errors.New("not configured: run `wl config`")
	}
	return cfg, nil
}

// Save writes the config with owner-only permissions and returns its path.
func Save(cfg Config) (string, error) {
	p, err := path()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return "", fmt.Errorf("create config dir: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encode config: %w", err)
	}
	if err := os.WriteFile(p, data, 0o600); err != nil {
		return "", fmt.Errorf("write %s: %w", p, err)
	}
	return p, nil
}
