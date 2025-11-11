package main

import (
	"fmt"
	"log"

	"github.com/DraconDev/go-templ-htmx-ex/db"
)

func main() {
	// Initialize database connection
	config := db.DefaultConfig()
	database, err := db.NewDatabase(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create tables
	if err := database.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Initialize repositories
	userRepo := db.NewUserRepository(database)
	prefsRepo := db.NewUserPreferencesRepository(database)

	// Example: Create a user
	user := &db.User{
		AuthID:  "auth123",
		Email:   "test@example.com",
		Name:    "Test User",
		Picture: "https://example.com/avatar.jpg",
	}

	if err := userRepo.CreateUser(user); err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("Created user with ID: %s\n", user.ID)

	// Example: Get user by auth ID
	foundUser, err := userRepo.GetUserByAuthID("auth123")
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	if foundUser != nil {
		fmt.Printf("Found user: %s (%s)\n", foundUser.Name, foundUser.Email)
	}

	// Example: Get user preferences
	prefs, err := prefsRepo.GetUserPreferences(user.ID)
	if err != nil {
		log.Fatalf("Failed to get user preferences: %v", err)
	}

	fmt.Printf("User preferences: Theme=%s, Language=%s\n", prefs.Theme, prefs.Language)

	fmt.Println("Database test completed successfully!")
}