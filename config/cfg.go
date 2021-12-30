package config

import (
	"os"
	"strconv"
	"strings"
)

type (
	Config struct {
		DatabaseFile string
		Host         string
		SessionKey   string
		CSRFKey      string
		CSSFile      string
		Title        string
		MOTDFile     string
		Topics       []string
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
		Topics:       strings.Split(os.Getenv("TOPICS"), ","),
	}
	cfg.Host = os.Getenv("HOST")
	if os.Getenv("HOST") == "" {
		cfg.Host = "localhost:8080"
	}
	perPage, _ := strconv.Atoi(os.Getenv("PER_PAGE"))
	cfg.PerPage = perPage
	if os.Getenv("PER_PAGE") == "" {
		cfg.PerPage = 50
	}
	return cfg
}
