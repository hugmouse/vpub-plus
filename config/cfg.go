package config

import (
	"os"
	"strconv"
)

type (
	Config struct {
		DatabaseURL string
		Port        string
		URL         string
		SessionKey  string
		CSRFKey     string
		CSSFile     string
		Title       string
		PerPage     int
	}
)

func New() *Config {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		SessionKey:  os.Getenv("SESSION_KEY"),
		CSRFKey:     os.Getenv("CSRF_KEY"),
		CSSFile:     os.Getenv("CSS_FILE"),
		Title:       os.Getenv("TITLE"),
		URL:         os.Getenv("URL"),
		Port:        os.Getenv("PORT"),
	}
	if os.Getenv("URL") == "" {
		cfg.URL = "http://localhost"
	}
	if os.Getenv("PORT") == "" {
		cfg.Port = "8080"
	}
	perPage, _ := strconv.Atoi(os.Getenv("PER_PAGE"))
	cfg.PerPage = perPage
	if os.Getenv("PER_PAGE") == "" {
		cfg.PerPage = 50
	}
	return cfg
}
