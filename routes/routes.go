package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/DraconDev/go-templ-htmx-ex/handlers"
	"github.com/DraconDev/go-templ-htmx-ex/handlers/admin"
	authHandlers "github.com/DraconDev/go-templ-htmx-ex/handlers/auth"
	"github.com/DraconDev/go-templ-htmx-ex/middleware"
)

// Route represents a route definition
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Middleware  []mux.MiddlewareFunc
}

// RouteDefinitions contains all application routes
var RouteDefinitions = []Route{
	// =============================================================================
	// PUBLIC ROUTES - No authentication required
	// =============================================================================

	// Homepage - Main landing page with platform showcase
	{
		Name:        "home",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: handlers.HomeHandler,
	},

	// Health check - API health monitoring endpoint
	{
		Name:        "health",
		Method:      "GET",
		Pattern:     "/health",
		HandlerFunc: handlers.HealthHandler,
	},

	// Login page - OAuth provider selection UI
	{
		Name:        "login",
		Method:      "GET",
		Pattern:     "/login",
		HandlerFunc: handlers.LoginHandler,
	},

	// =============================================================================
	// OAUTH AUTHENTICATION FLOW
	// =============================================================================

	// OAuth Login Route - Consolidated with provider parameter
	{
		Name:        "oauth_login",
		Method:      "GET",
		Pattern:     "/auth/login",
		HandlerFunc: authHandlers.LoginHandler,
	},

	// OAuth callback handler
	{
		Name:        "oauth_callback",
		Method:      "GET",
		Pattern:     "/auth/callback",
		HandlerFunc: authHandlers.AuthCallbackHandler,
	},

	// =============================================================================
	// PROTECTED USER ROUTES - Authentication required
	// =============================================================================

	// User profile page - Display user information and account details
	{
		Name:        "profile",
		Method:      "GET",
		Pattern:     "/profile",
		HandlerFunc: handlers.ProfileHandler,
	},

	// =============================================================================
	// ADMIN ROUTES - Admin authentication required
	// =============================================================================

	// Admin dashboard - Main admin interface for platform management
	{
		Name:        "admin_dashboard",
		Method:      "GET",
		Pattern:     "/admin",
		HandlerFunc: admin.AdminDashboardHandler,
	},

	// Admin API endpoints - Management operations for administrators
	{
		Name:        "admin_get_users",
		Method:      "GET",
		Pattern:     "/api/admin/users",
		HandlerFunc: admin.GetUsersHandler,
	},

	{
		Name:        "admin_get_analytics",
		Method:      "GET",
		Pattern:     "/api/admin/analytics",
		HandlerFunc: admin.GetAnalyticsHandler,
	},

	{
		Name:        "admin_get_settings",
		Method:      "GET",
		Pattern:     "/api/admin/settings",
		HandlerFunc: admin.GetSettingsHandler,
	},

	{
		Name:        "admin_get_logs",
		Method:      "GET",
		Pattern:     "/api/admin/logs",
		HandlerFunc: admin.GetLogsHandler,
	},

	// =============================================================================
	// SESSION MANAGEMENT API - Authentication required
	// =============================================================================

	// Logout user - Destroy current session and clear cookies
	{
		Name:        "logout",
		Method:      "POST",
		Pattern:     "/api/auth/logout",
		HandlerFunc: authHandlers.LogoutHandler,
	},

	// Set session - Create new server session with provided session ID
	{
		Name:        "set_session",
		Method:      "POST",
		Pattern:     "/api/auth/set-session",
		HandlerFunc: authHandlers.SetSessionHandler,
	},
}

// SetupRoutes configures and returns the router with all routes
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Add authentication middleware to all routes
	router.Use(middleware.AuthMiddleware)

	// Register all routes
	for _, route := range RouteDefinitions {
		router.HandleFunc(route.Pattern, route.HandlerFunc).
			Methods(route.Method).
			Name(route.Name)
	}

	// Static files (for CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return router
}