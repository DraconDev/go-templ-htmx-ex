package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Create authentication service
	authService := auth.NewService(cfg.AuthServiceURL)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService, cfg)

	// Create router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	router.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")

	// OAuth login routes
	fmt.Printf("🔗 REGISTERING: /auth/google (GET)\n")
	router.HandleFunc("/auth/google", authHandler.GoogleLoginHandler).Methods("GET")
	fmt.Printf("🔗 REGISTERED: /auth/google\n")
	
	fmt.Printf("🔗 REGISTERING: /auth/github (GET)\n")
	router.HandleFunc("/auth/github", authHandler.GitHubLoginHandler).Methods("GET")
	fmt.Printf("🔗 REGISTERED: /auth/github\n")
	
	// Auth service callback endpoints
	fmt.Printf("🔗 REGISTERING: /auth/google/callback (GET)\n")
	router.HandleFunc("/auth/google/callback", authHandler.AuthCallbackHandler).Methods("GET")
	fmt.Printf("🔗 REGISTERED: /auth/google/callback\n")
	
	fmt.Printf("🔗 REGISTERING: /auth/github/callback (GET)\n")
	router.HandleFunc("/auth/github/callback", authHandler.AuthCallbackHandler).Methods("GET")
	fmt.Printf("🔗 REGISTERED: /auth/github/callback\n")
	
	fmt.Printf("🔗 REGISTERING: /auth/callback (GET)\n")
	router.HandleFunc("/auth/callback", authHandler.AuthCallbackHandler).Methods("GET")
	fmt.Printf("🔗 REGISTERED: /auth/callback\n")

	// Session management API
	router.HandleFunc("/api/auth/validate", authHandler.ValidateSessionHandler).Methods("POST")
	router.HandleFunc("/api/auth/logout", authHandler.LogoutHandler).Methods("POST")
	router.HandleFunc("/api/auth/user", authHandler.GetUserHandler).Methods("GET")
	router.HandleFunc("/api/auth/set-session", authHandler.SetSessionHandler).Methods("POST")

	// Static files (for CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.ServerPort)
		log.Printf("Visit http://localhost:%s to access the application", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
