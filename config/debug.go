package config

import (
	"log"
	"os"
	"strings"
)

// DebugConfig holds debugging settings
type DebugConfig struct {
	Enabled        bool
	Auth           bool
	Requests       bool
	Templates      bool
	Token          bool
	OAuth          bool
}

// LoadDebugConfig loads debug configuration from environment
func LoadDebugConfig() *DebugConfig {
	debug := os.Getenv("DEBUG")
	if debug == "" {
		debug = os.Getenv("ENV")
		if debug == "" {
			debug = "false"
		}
	}

	return &DebugConfig{
		Enabled:        strings.ToLower(debug) == "debug" || strings.ToLower(debug) == "true",
		Auth:           true, // Always enable auth debugging
		Requests:       strings.Contains(strings.ToLower(debug), "request"),
		Templates:      strings.Contains(strings.ToLower(debug), "template"),
		Token:          strings.Contains(strings.ToLower(debug), "token"),
		OAuth:          strings.Contains(strings.ToLower(debug), "oauth"),
	}
}

// LogAuth logs authentication-related messages
func (d *DebugConfig) LogAuth(msg string, args ...interface{}) {
	if d.Enabled && d.Auth {
		log.Printf("🔐 AUTH: "+msg, args...)
	}
}

// LogRequest logs request-related messages
func (d *DebugConfig) LogRequest(msg string, args ...interface{}) {
	if d.Enabled && d.Requests {
		log.Printf("📡 REQUEST: "+msg, args...)
	}
}

// LogTemplate logs template-related messages
func (d *DebugConfig) LogTemplate(msg string, args ...interface{}) {
	if d.Enabled && d.Templates {
		log.Printf("🎨 TEMPLATE: "+msg, args...)
	}
}

// LogToken logs token-related messages
func (d *DebugConfig) LogToken(msg string, args ...interface{}) {
	if d.Enabled && d.Token {
		log.Printf("🔑 TOKEN: "+msg, args...)
	}
}

// LogOAuth logs OAuth-related messages
func (d *DebugConfig) LogOAuth(msg string, args ...interface{}) {
	if d.Enabled && d.OAuth {
		log.Printf("🔓 OAUTH: "+msg, args...)
	}
}