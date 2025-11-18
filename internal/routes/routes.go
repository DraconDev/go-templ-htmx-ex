package routes

import (
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/internal/handlers"
	"github.com/DraconDev/go-templ-htmx-ex/internal/handlers/admin"
	"github.com/DraconDev/go-templ-htmx-ex/internal/handlers/auth/login"
	"github.com/DraconDev/go-templ-htmx-ex/internal/handlers/auth/session"
	"github.com/DraconDev/go-templ-htmx-ex/internal/middleware"
	"github.com/gorilla/mux"
)

// HandlerInstances holds all handler instances for route registration
type HandlerInstances struct {
	AdminHandler   *admin.AdminHandler
	LoginHandler   *login.LoginHandler
	SessionHandler *session.SessionHandler
}

// SetupRoutes configures and returns the router with all routes
// This approach accepts handler instances to avoid method reference issues
func SetupRoutes(handlerInstances *HandlerInstances) *mux.Router {
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
	if handlerInstances.LoginHandler != nil {
		router.HandleFunc("/auth/login", handlerInstances.LoginHandler.LoginHandler).Methods("GET")
		router.HandleFunc("/auth/callback", handlerInstances.LoginHandler.AuthCallbackHandler).Methods("GET")
	}

	// =============================================================================
	// PROTECTED USER ROUTES - Authentication required
	// =============================================================================

	// User profile page - Display user information and account details
	router.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")

	// =============================================================================
	// ADMIN ROUTES - Admin authentication required
	// =============================================================================

	// Admin dashboard - Main admin interface for platform management
	if handlerInstances.AdminHandler != nil {
		router.HandleFunc("/admin", handlerInstances.AdminHandler.AdminDashboardHandler).Methods("GET")
		router.HandleFunc("/api/admin/users", handlerInstances.AdminHandler.GetUsersHandler).Methods("GET")
		router.HandleFunc("/api/admin/analytics", handlerInstances.AdminHandler.GetAnalyticsHandler).Methods("GET")
		router.HandleFunc("/api/admin/settings", handlerInstances.AdminHandler.GetSettingsHandler).Methods("GET")
		router.HandleFunc("/api/admin/logs", handlerInstances.AdminHandler.GetLogsHandler).Methods("GET")
	}

	// =============================================================================
	// SESSION MANAGEMENT API - Authentication required
	// =============================================================================

	if handlerInstances.SessionHandler != nil {
		// Logout user - Destroy current session and clear cookies
		router.HandleFunc("/api/auth/logout", handlerInstances.SessionHandler.LogoutHandler).Methods("POST")

		// Set session - Create new server session with provided session ID
		router.HandleFunc("/api/auth/set-session", handlerInstances.SessionHandler.SetSessionHandler).Methods("POST")

		// Exchange code - Exchange OAuth authorization code for session tokens
		router.HandleFunc("/api/auth/exchange-code", handlerInstances.SessionHandler.ExchangeCodeHandler).Methods("POST")
	}

	// Static files (for CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return router
}
