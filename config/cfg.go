// Environment Variables:
//   - DATABASE_URL: PostgreSQL connection string
//   - PORT: HTTP server port (default: "8080")
//   - SESSION_KEY: 32-byte session encryption key
//   - CSRF_KEY: 32-byte CSRF protection key
//   - TITLE: Forum title (default: "My vpub-plus forum")
//   - CSRF_SECURE: Enable secure cookies (default: true)
//   - PROXYING_ENABLED: Enable image proxying (default: true)
//   - POSTGRES_MAX_OPEN_CONNECTIONS: Max DB connections (default: 0 - auto)
//   - POSTGRES_MAX_IDLE_CONNECTIONS: Idle DB connections (default: 0 - auto)
//   - POSTGRES_MAX_LIFETIME: Connection max lifetime (default: 5m)
package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL                string        // PostgreSQL connection string
	DatabaseMaxOpenConnections int           // Maximum number of open DB connections
	DatabaseMaxIdleConnections int           // Maximum number of idle DB connections
	DatabaseMaxLifetime        time.Duration // Maximum lifetime of DB connections
	Port                       string        // HTTP server port
	SessionKey                 string        // Session encryption key (32 bytes)
	CSRFKey                    string        // CSRF protection key (32 bytes)
	CSRFSecure                 bool          // Use secure cookies
	Title                      string        // Forum title
	ProxyingEnabled            bool          // Enable image proxying
}

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
