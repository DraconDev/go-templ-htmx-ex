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

	// Create admin handler with SQLC queries (handle nil db gracefully)
	if queries != nil {
		adminHandler = handlers.NewAdminHandler(cfg, queries)
	} else {
		log.Println("⚠️  Admin handler not initialized - no database connection")
	}

	// Create auth handler
	authHandler = handlers.NewAuthHandler(cfg)

	// Create router
	router := mux.NewRouter()

	// Add authentication middleware to all routes
	router.Use(middleware.AuthMiddleware)

	// =============================================================================
	// PUBLIC ROUTES - No authentication required
	// =============================================================================

	// Homepage - Main landing page with platform showcase
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Health check - API health monitoring endpoint
	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")

	// Login page - OAuth provider selection UI
	router.HandleFunc("/login", handlers.LoginHandler).Methods("GET")

	// =============================================================================
	// OAUTH AUTHENTICATION FLOW
	// =============================================================================

	// Step 1: Redirect to OAuth providers
	// OAuth Login Route - Consolidated with provider parameter
	router.HandleFunc("/auth/login", authHandler.LoginHandler).Methods("GET")

	// Step 2: OAuth callback handler
	// This route receives the OAuth callback, processes tokens, and creates server sessions
	router.HandleFunc("/auth/callback", authHandler.AuthCallbackHandler).Methods("GET")

	// =============================================================================
	// DEVELOPMENT & TESTING ROUTES
	// =============================================================================

	// Test authentication page - Development tool for testing OAuth flows
	// This route is now handled by the auth_test.go file
	// router.HandleFunc("/test", authHandler.TestTokenRefreshHandler).Methods("GET")

	// =============================================================================
	// PROTECTED USER ROUTES - Authentication required
	// =============================================================================

	// User profile page - Display user information and account details
	router.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")

	// =============================================================================
	// ADMIN ROUTES - Admin authentication required
	// =============================================================================

	// Admin dashboard - Main admin interface for platform management
	router.HandleFunc("/admin", adminHandler.AdminDashboardHandler).Methods("GET")

	// Admin API endpoints - Management operations for administrators
	router.HandleFunc("/api/admin/users", adminHandler.GetUsersHandler).Methods("GET")         // List all users
	router.HandleFunc("/api/admin/analytics", adminHandler.GetAnalyticsHandler).Methods("GET") // Platform analytics
	router.HandleFunc("/api/admin/settings", adminHandler.GetSettingsHandler).Methods("GET")   // System settings
	router.HandleFunc("/api/admin/logs", adminHandler.GetLogsHandler).Methods("GET")           // System logs

	// =============================================================================
	// SESSION MANAGEMENT API - Authentication required
	// =============================================================================

	

	// Logout user - Destroy current session and clear cookies
	router.HandleFunc("/api/auth/logout", authHandler.LogoutHandler).Methods("POST")

	

	// Set session - Create new server session with provided session ID
	router.HandleFunc("/api/auth/set-session", authHandler.SetSessionHandler).Methods("POST")

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
