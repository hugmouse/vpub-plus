package config

import (
	"os"
)

type (
	Config struct {
		DatabaseURL     string
		Port            string
		SessionKey      string
		CSRFKey         string
		CSRFSecure      bool
		Title           string
		ProxyingEnabled bool
	}
)

func New() *Config {
	cfg := &Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		SessionKey:      os.Getenv("SESSION_KEY"),
		CSRFKey:         os.Getenv("CSRF_KEY"),
		CSRFSecure:      os.Getenv("CSRF_SECURE") == "true",
		Title:           os.Getenv("TITLE"),
		Port:            os.Getenv("PORT"),
		ProxyingEnabled: os.Getenv("PROXYING_ENABLED") == "true",
	}
	if os.Getenv("PORT") == "" {
		cfg.Port = "8080"
	}
	return cfg
}
