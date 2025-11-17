package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/DraconDev/go-templ-htmx-ex/handlers"
	"github.com/DraconDev/go-templ-htmx-ex/handlers/admin"
	authHandlers "github.com/DraconDev/go-templ-htmx-ex/handlers/auth"
	"github.com/DraconDev/go-templ-htmx-ex/middleware"
)

// HandlerInstances holds all handler instances for route registration
type HandlerInstances struct {
	AuthHandler *authHandlers.AuthHandler
	AdminHandler *admin.AdminHandler
}

// SetupRoutes configures and returns the router with all routes
// This approach accepts handler instances to avoid method reference issues
func SetupRoutes(handlers *HandlerInstances) *mux.Router {
	router := mux.NewRouter()

	// Add authentication middleware to all routes
	router.Use(middleware.AuthMiddleware)

	// =============================================================================
	// PUBLIC ROUTES - No authentication required
	// =============================================================================

	// Homepage - Main landing page with platform showcase
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET").Name("home")

	// Health check - API health monitoring endpoint
	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET").Name("health")

	// Login page - OAuth provider selection UI
	router.HandleFunc("/login", handlers.LoginHandler).Methods("GET").Name("login")

	// =============================================================================
	// OAUTH AUTHENTICATION FLOW
	// =============================================================================

	// OAuth Login Route - Consolidated with provider parameter
	router.HandleFunc("/auth/login", handlers.AuthHandler.LoginHandler).Methods("GET").Name("oauth_login")

	// OAuth callback handler
	router.HandleFunc("/auth/callback", handlers.AuthHandler.AuthCallbackHandler).Methods("GET").Name("oauth_callback")

	// =============================================================================
	// PROTECTED USER ROUTES - Authentication required
	// =============================================================================

	// User profile page - Display user information and account details
	router.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET").Name("profile")

	// =============================================================================
	// ADMIN ROUTES - Admin authentication required
	// =============================================================================

	// Admin dashboard - Main admin interface for platform management
	if handlers.AdminHandler != nil {
		router.HandleFunc("/admin", handlers.AdminHandler.AdminDashboardHandler).Methods("GET").Name("admin_dashboard")
		router.HandleFunc("/api/admin/users", handlers.AdminHandler.GetUsersHandler).Methods("GET").Name("admin_get_users")
		router.HandleFunc("/api/admin/analytics", handlers.AdminHandler.GetAnalyticsHandler).Methods("GET").Name("admin_get_analytics")
		router.HandleFunc("/api/admin/settings", handlers.AdminHandler.GetSettingsHandler).Methods("GET").Name("admin_get_settings")
		router.HandleFunc("/api/admin/logs", handlers.AdminHandler.GetLogsHandler).Methods("GET").Name("admin_get_logs")
	}

	// =============================================================================
	// SESSION MANAGEMENT API - Authentication required
	// =============================================================================

	// Logout user - Destroy current session and clear cookies
	router.HandleFunc("/api/auth/logout", handlers.AuthHandler.LogoutHandler).Methods("POST").Name("logout")

	// Set session - Create new server session with provided session ID
	router.HandleFunc("/api/auth/set-session", handlers.AuthHandler.SetSessionHandler).Methods("POST").Name("set_session")

	// =============================================================================
	// STATIC FILES
	// =============================================================================

	// Static files (for CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return router
}

// RouteSummary provides a summary of all registered routes
type RouteSummary struct {
	TotalRoutes     int `json:"total_routes"`
	PublicRoutes    int `json:"public_routes"`
	ProtectedRoutes int `json:"protected_routes"`
	AdminRoutes     int `json:"admin_routes"`
	AuthAPIRoutes   int `json:"auth_api_routes"`
}

// CountRoutes provides a count of all route types
func CountRoutes() RouteSummary {
	return RouteSummary{
		TotalRoutes:     15, // Approximate count
		PublicRoutes:    3,
		ProtectedRoutes: 1,
		AdminRoutes:     5,
		AuthAPIRoutes:   4,
	}
}