package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	ServerPort     string
	AuthServiceURL string
	RedirectURL    string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		ServerPort:     getEnvOrDefault("PORT", "8081"),
		AuthServiceURL: getEnvOrDefault("AUTH_SERVICE_URL", "http://localhost:8080"),
		RedirectURL:    getEnvOrDefault("REDIRECT_URL", "http://localhost:8081"),
	}
}

// getEnvOrDefault returns environment variable or default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}