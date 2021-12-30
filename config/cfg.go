package config

import "os"

type (
	Config struct {
		DatabaseFile string
		Host         string
		SessionKey   string
		CSRFKey      string
		CSSFile      string
		Title        string
		MOTDFile     string
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
	}
	cfg.Host = os.Getenv("HOST")
	if os.Getenv("HOST") == "" {
		cfg.Host = "localhost:8080"
	}
	return cfg
}
