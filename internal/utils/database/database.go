package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// DB holds the database connection pool
var DB *sql.DB

// InitDatabase initializes the database connection using DB_URL environment variable
func InitDatabase() error {
	// Get database URL from environment
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return fmt.Errorf("DB_URL environment variable is required")
	}

	// Open database connection
	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connection established successfully with Neon DB")
	return nil
}

// InitDatabaseFromConnString initializes database from a connection string
func InitDatabaseFromConnString(connString string) error {
	var err error
	DB, err = sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connection established successfully")
	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB returns the database connection pool
func GetDB() *sql.DB {
	return DB
}

// IsInitialized checks if the database connection pool is initialized
func IsInitialized() bool {
	return DB != nil
}