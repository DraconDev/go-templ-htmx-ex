package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Database represents the database connection and operations
type Database struct {
	*sql.DB
}

// Config represents database configuration
type Config struct {
	URL      string
	MaxOpen  int
	MaxIdle  int
	MaxLife  time.Duration
}

// DefaultConfig returns default database configuration
func DefaultConfig() *Config {
	return &Config{
		URL:     getEnvOrDefault("DB_URL", ""),
		MaxOpen: 25,
		MaxIdle: 10,
		MaxLife: 5 * time.Minute,
	}
}

// getEnvOrDefault returns environment variable or default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewDatabase creates a new database connection
func NewDatabase(config *Config) (*Database, error) {
	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, config.SSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpen)
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetConnMaxLifetime(config.MaxLife)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db}, nil
}

// CreateTables creates the database schema
func (db *Database) CreateTables() error {
	schema := `
	-- Users table: Application-level user data
	-- Complements external auth service
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		auth_id VARCHAR(255) UNIQUE NOT NULL, -- References auth service user ID
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		picture TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- User preferences: Essential user settings
	CREATE TABLE IF NOT EXISTS user_preferences (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		theme VARCHAR(20) DEFAULT 'dark',
		language VARCHAR(10) DEFAULT 'en',
		email_notifications BOOLEAN DEFAULT true,
		push_notifications BOOLEAN DEFAULT true,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_users_auth_id ON users(auth_id);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);

	-- Function to update updated_at timestamp
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
	END;
	$$ language 'plpgsql';

	-- Triggers for updated_at
	CREATE TRIGGER update_users_updated_at 
		BEFORE UPDATE ON users
		FOR EACH ROW 
		EXECUTE FUNCTION update_updated_at_column();

	CREATE TRIGGER update_user_preferences_updated_at 
		BEFORE UPDATE ON user_preferences
		FOR EACH ROW 
		EXECUTE FUNCTION update_updated_at_column();
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database tables created successfully")
	return nil
}

// HealthCheck performs a database health check
func (db *Database) HealthCheck() error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.DB.Close()
}