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

	dbSqlc "github.com/DraconDev/go-templ-htmx-ex/database/sqlc"
	"github.com/DraconDev/go-templ-htmx-ex/internal/handlers"
	"github.com/DraconDev/go-templ-htmx-ex/internal/handlers/admin"
	"github.com/DraconDev/go-templ-htmx-ex/internal/middleware"
	"github.com/DraconDev/go-templ-htmx-ex/internal/services"
	"github.com/DraconDev/go-templ-htmx-ex/internal/utils/config"
	database "github.com/DraconDev/go-templ-htmx-ex/internal/utils/database"
	_ "github.com/lib/pq"
)

var adminHandler *admin.AdminHandler
var userService *services.UserService
var authService *services.AuthService
var sqlDB *sql.DB
var queries *dbSqlc.Queries

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database if configured
	if err := database.InitDatabaseIfConfigured(); err != nil {
		log.Printf("⚠️  Database initialization failed: %v", err)
		log.Println("💡 Continuing without database functionality")
	}

	// Get database URL from environment for runtime connection
	dbURL := os.Getenv("DB_URL")
	if dbURL != "" {
		log.Printf("🔗 Connecting to database for runtime...")
		var err error
		sqlDB, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Printf("❌ Database connection failed: %v", err)
			log.Println("⚠️  Continuing without database...")
			sqlDB = nil
		} else {
			// Test connection
			if err := sqlDB.Ping(); err != nil {
				log.Printf("❌ Database ping failed: %v", err)
				log.Println("⚠️  Continuing without database...")
				sqlDB = nil
			} else {
				log.Println("✅ Database connected successfully")

				// Initialize SQLC queries
				queries = dbSqlc.New(sqlDB)
				log.Println("✅ SQLC queries initialized")
			}
		}
	}

	// Initialize services
	if queries != nil {
		userService = services.NewUserService(queries)
		log.Println("✅ User service initialized")
	}

	authService = services.NewAuthService(cfg)
	log.Println("✅ Auth service initialized")

	// Create handlers with services
	if queries != nil {
		adminHandler = admin.NewAdminHandler(cfg, queries)
	} else {
		log.Println("⚠️  Admin handler not initialized - no database connection")
	}

	loginHandler := handlers.LoginHandler.NewAuthHandler(cfg)

	// Create router using new route structure
	router := SetupRoutes()

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

// SetupRoutes creates and configures the router with all routes
func SetupRoutes() *mux.Router {
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

	// OAuth Login Route - Consolidated with provider parameter
	router.HandleFunc("/auth/login", handlers.LoginHandler).Methods("GET")

	// OAuth callback handler
	router.HandleFunc("/auth/callback", handlers.AuthCallbackHandler).Methods("GET")

	// =============================================================================
	// PROTECTED USER ROUTES - Authentication required
	// =============================================================================

	// User profile page - Display user information and account details
	router.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")

	// =============================================================================
	// ADMIN ROUTES - Admin authentication required
	// =============================================================================

	// Admin dashboard - Main admin interface for platform management
	if adminHandler != nil {
		router.HandleFunc("/admin", adminHandler.AdminDashboardHandler).Methods("GET")
		router.HandleFunc("/api/admin/users", adminHandler.GetUsersHandler).Methods("GET")
		router.HandleFunc("/api/admin/analytics", adminHandler.GetAnalyticsHandler).Methods("GET")
		router.HandleFunc("/api/admin/settings", adminHandler.GetSettingsHandler).Methods("GET")
		router.HandleFunc("/api/admin/logs", adminHandler.GetLogsHandler).Methods("GET")
	}

	// =============================================================================
	// SESSION MANAGEMENT API - Authentication required
	// =============================================================================

	// Logout user - Destroy current session and clear cookies
	router.HandleFunc("/api/auth/logout", authHandler.LogoutHandler).Methods("POST")

	// Set session - Create new server session with provided session ID
	router.HandleFunc("/api/auth/set-session", authHandler.SetSessionHandler).Methods("POST")

	// Static files (for CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return router
}
