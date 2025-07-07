package config

import (
	"os"
	"strconv"
	"time" // Import time for duration parsing
)

type (
	Config struct {
		DatabaseURL                string
		DatabaseMaxOpenConnections int
		DatabaseMaxIdleConnections int
		DatabaseMaxLifetime        time.Duration
		Port                       string
		SessionKey                 string
		CSRFKey                    string
		CSRFSecure                 bool
		Title                      string
		ProxyingEnabled            bool
	}
)

func New() *Config {
	cfg := &Config{
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		Port:            getEnv("PORT", "8080"),
		SessionKey:      getEnv("SESSION_KEY", "your32byteslongsessionkeyhere"),
		CSRFKey:         getEnv("CSRF_KEY", "your32byteslongcsrfkeyhere"),
		Title:           getEnv("TITLE", "My vpub-plus forum"),
		CSRFSecure:      getEnvAsBool("CSRF_SECURE", true),
		ProxyingEnabled: getEnvAsBool("PROXYING_ENABLED", true),
	}

	cfg.DatabaseMaxOpenConnections = getEnvAsInt("POSTGRES_MAX_OPEN_CONNECTIONS", 0)
	cfg.DatabaseMaxIdleConnections = getEnvAsInt("POSTGRES_MAX_IDLE_CONNECTIONS", 0)
	cfg.DatabaseMaxLifetime = getEnvAsDuration("POSTGRES_MAX_LIFETIME", 5*time.Minute)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
