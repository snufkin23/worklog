package platform

import (
	"errors"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	APIToken    string
	CronSecret  string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		Port:        getenv("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		APIToken:    os.Getenv("API_TOKEN"),
		CronSecret:  os.Getenv("CRON_SECRET"),
	}
	if cfg.DatabaseURL == "" || cfg.APIToken == "" || cfg.CronSecret == "" {
		return Config{}, errors.New("DATABASE_URL, API_TOKEN and CRON_SECRET must be set")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
