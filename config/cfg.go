package config

import "os"

type (
	Config struct {
		DatabaseFile string
		Host         string
		SessionKey   string
		Env          string
		CSRFKey      string
	}
)

func New() *Config {
	cfg := &Config{
		DatabaseFile: os.Getenv("DATABASE_FILE"),
		SessionKey:   os.Getenv("SESSION_KEY"),
		Env:          os.Getenv("ENV"),
		CSRFKey:      os.Getenv("CSRF_KEY"),
	}
	cfg.Host = os.Getenv("HOST")
	if os.Getenv("HOST") == "" {
		cfg.Host = "localhost:8080"
	}
	return cfg
}
