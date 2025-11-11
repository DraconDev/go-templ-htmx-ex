package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	ServerPort     string
	AuthServiceURL string
	RedirectURL    string
	AuthSecret     string
}

var (
	// Global config instance
	Current *Config
)

// getEnvOrDefault returns environment variable or default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	config := &Config{
		ServerPort:     getEnvOrDefault("PORT", "8081"),
		AuthServiceURL: getEnvOrDefault("AUTH_SERVICE_URL", "http://localhost:8080"),
		RedirectURL:    getEnvOrDefault("REDIRECT_URL", "http://localhost:8081"),
		AuthSecret:     getEnvOrDefault("AUTH_SECRET", ""),
	}

	Current = config
	return config
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%s", c.ServerPort)
}