package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/DraconDev/go-templ-htmx-ex/auth"
	"github.com/DraconDev/go-templ-htmx-ex/config"
	dbInit "github.com/DraconDev/go-templ-htmx-ex/db"
	dbSqlc "github.com/DraconDev/go-templ-htmx-ex/db/sqlc"
	"github.com/DraconDev/go-templ-htmx-ex/handlers"
	"github.com/DraconDev/go-templ-htmx-ex/middleware"
	_ "github.com/lib/pq"
)

var authHandler *handlers.AuthHandler
var adminHandler *handlers.AdminHandler
var db *sql.DB
var queries *dbSqlc.Queries

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database if configured
	if err := dbInit.InitDatabaseIfConfigured(); err != nil {
		log.Printf("⚠️  Database initialization failed: %v", err)
		log.Println("💡 Continuing without database functionality")
	}

	// Get database URL from environment for runtime connection
	dbURL := os.Getenv("DB_URL")
	if dbURL != "" {
		log.Printf("🔗 Connecting to database for runtime...")
		var err error
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Printf("❌ Database connection failed: %v", err)
			log.Println("⚠️  Continuing without database...")
			db = nil
		} else {
			// Test connection
			if err := db.Ping(); err != nil {
				log.Printf("❌ Database ping failed: %v", err)
				log.Println("⚠️  Continuing without database...")
				db = nil
			} else {
				log.Println("✅ Database connected successfully")

				// Initialize SQLC queries
				queries = dbSqlc.New(db)
				log.Println("✅ SQLC queries initialized")
			}
		}
	}

	// Create auth service
	authService := auth.NewService(cfg)

	// Create admin handler with SQLC queries (handle nil db gracefully)
	if queries != nil {
		adminHandler = handlers.NewAdminHandler(cfg, queries)
	} else {
		log.Println("⚠️  Admin handler not initialized - no database connection")
	}

	// Create auth handler
	authHandler = handlers.NewAuthHandler(authService, cfg)

	// Create router
	router := mux.NewRouter()

	// Add authentication middleware to all routes
	router.Use(middleware.AuthMiddleware)

	// Define routes using clean handlers from handlers/handlers.go
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("GET")

	// OAuth login routes
	router.HandleFunc("/auth/google", authHandler.GoogleLoginHandler).Methods("GET")
	router.HandleFunc("/auth/github", authHandler.GitHubLoginHandler).Methods("GET")
	router.HandleFunc("/auth/discord", authHandler.DiscordLoginHandler).Methods("GET")
	router.HandleFunc("/auth/microsoft", authHandler.MicrosoftLoginHandler).Methods("GET")
	router.HandleFunc("/auth/callback", authHandler.AuthCallbackHandler).Methods("GET")

	// Test routes
	router.HandleFunc("/test", authHandler.TestTokenRefreshHandler).Methods("GET")

	// User profile page
	router.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")

	// Admin dashboard
	router.HandleFunc("/admin", adminHandler.AdminDashboardHandler).Methods("GET")

	// Admin API routes
	router.HandleFunc("/api/admin/users", adminHandler.GetUsersHandler).Methods("GET")
	router.HandleFunc("/api/admin/analytics", adminHandler.GetAnalyticsHandler).Methods("GET")
	router.HandleFunc("/api/admin/settings", adminHandler.GetSettingsHandler).Methods("GET")
	router.HandleFunc("/api/admin/logs", adminHandler.GetLogsHandler).Methods("GET")

	// Session management
	router.HandleFunc("/api/auth/validate", authHandler.ValidateSessionHandler).Methods("POST")
	router.HandleFunc("/api/auth/logout", authHandler.LogoutHandler).Methods("POST")
	router.HandleFunc("/api/auth/refresh", authHandler.RefreshTokenHandler).Methods("POST")
	router.HandleFunc("/api/auth/user", authHandler.GetUserHandler).Methods("GET")
	router.HandleFunc("/api/auth/set-session", authHandler.SetSessionHandler).Methods("POST")
	router.HandleFunc("/api/auth/exchange-code", authHandler.ExchangeCodeHandler).Methods("POST")

	// Static files (for CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
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
