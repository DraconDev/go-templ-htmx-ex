package config

import (
	"fmt"
	"log"

	"github.com/dracondev/go-templ-htmx-ex/libs/configx"
)

// Config holds application configuration
type Config struct {
	*configx.Config
	ServerPort     string
	AuthServiceURL string
	RedirectURL    string
	AdminEmail     string
}

var (
	// Global config instance
	Current *Config
)

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	fields := []configx.ConfigField{
		{
			Key:          "PORT",
			DefaultValue: "8081",
			Required:     false,
			Description:  "Server port",
		},
		{
			Key:          "AUTH_SERVICE_URL",
			DefaultValue: "http://localhost:8080",
			Required:     false,
			Description:  "Authentication service URL",
		},
		{
			Key:          "REDIRECT_URL",
			DefaultValue: "http://localhost:8081",
			Required:     false,
			Description:  "OAuth redirect URL",
		},
		{
			Key:          "ADMIN_EMAIL",
			DefaultValue: "admin@startup-platform.local",
			Required:     false,
			Description:  "Admin email address",
		},
	}

	baseConfig, err := configx.Load(fields, configx.DefaultOptions())
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	config := &Config{
		Config:         baseConfig,
		ServerPort:     baseConfig.Get("PORT"),
		AuthServiceURL: baseConfig.Get("AUTH_SERVICE_URL"),
		RedirectURL:    baseConfig.Get("REDIRECT_URL"),
		AdminEmail:     baseConfig.Get("ADMIN_EMAIL"),
	}

	Current = config
	return config
}

// IsAdmin checks if the given email matches the admin email
func (c *Config) IsAdmin(email string) bool {
	return c.AdminEmail != "" && email == c.AdminEmail
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%s", c.ServerPort)
}
