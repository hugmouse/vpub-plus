package config

import (
	"os"
	"strconv"
)

type (
	Config struct {
		DatabaseFile string
		Port         string
		URL          string
		SessionKey   string
		CSRFKey      string
		CSSFile      string
		Title        string
		MOTDFile     string
		PerPage      int
	}
)

func New() *Config {
	cfg := &Config{
		DatabaseFile: os.Getenv("DATABASE_FILE"),
		SessionKey:   os.Getenv("SESSION_KEY"),
		CSRFKey:      os.Getenv("CSRF_KEY"),
		CSSFile:      os.Getenv("CSS_FILE"),
		Title:        os.Getenv("TITLE"),
		MOTDFile:     os.Getenv("MOTD_FILE"),
		URL:          os.Getenv("URL"),
		Port:         os.Getenv("PORT"),
	}
	if os.Getenv("URL") == "" {
		cfg.URL = "http://localhost"
	}
	if os.Getenv("PORT") == "" {
		cfg.Port = "8080"
	}
	if os.Getenv("TITLE") == "" {
		cfg.Title = "vpub"
	}
	if os.Getenv("DATABASE_FILE") == "" {
		cfg.DatabaseFile = "./vpub.sqlite"
	}
	perPage, _ := strconv.Atoi(os.Getenv("PER_PAGE"))
	cfg.PerPage = perPage
	if os.Getenv("PER_PAGE") == "" {
		cfg.PerPage = 50
	}
	return cfg
}
