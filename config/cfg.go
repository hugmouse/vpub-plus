package config

import (
	"os"
)

type (
	Config struct {
		DatabaseURL string
		Port        string
		SessionKey  string
		CSRFKey     string
		Title       string
		Address     string
	}
)

func New() *Config {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		SessionKey:  os.Getenv("SESSION_KEY"),
		CSRFKey:     os.Getenv("CSRF_KEY"),
		Title:       os.Getenv("TITLE"),
		Port:        os.Getenv("PORT"),
		Address:     os.Getenv("ADDRESS"),
	}
	if os.Getenv("PORT") == "" {
		cfg.Port = "8080"
	}
	return cfg
}
