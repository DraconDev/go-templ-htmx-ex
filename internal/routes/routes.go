package routes

import (
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/internal/middleware"
	"github.com/gorilla/mux"
)

// HandlerInstances holds all handler instances for route registration
type HandlerInstances struct {
	// These will be populated with actual handler instances
}

// Route categories for middleware configuration
// This allows easy configuration of which routes require authentication

// Protected routes that require authentication
var ProtectedRoutes = []string{
	"/profile", // User profile page
	"/admin",   // Admin dashboard
}

// Admin routes that require authentication
var AdminRoutes = []string{
	"/api/admin", // All admin API endpoints
}

// Public routes that don't require authentication
var PublicRoutes = []string{
	"/",       // Homepage
	"/health", // Health check endpoint
	"/login",  // Login page
	"/test",   // Test page for development
}

// OAuth routes that are part of authentication flow
var OAuthRoutes = []string{
	"/auth/google",
	"/auth/github",
	"/auth/discord",
	"/auth/microsoft",
	"/auth/callback",
}

// Auth API routes that handle session management
var AuthAPIRoutes = []string{
	"/api/auth/validate",
	"/api/auth/user",
	"/api/auth/logout",
	"/api/auth/set-session",
	"/api/auth/exchange-code",
	"/api/auth/test-session-create",
	"/api/auth/refresh",
}

// All routes - useful for validation
var AllRoutes = append(
	append(ProtectedRoutes, AdminRoutes...),
	append(PublicRoutes, append(OAuthRoutes, AuthAPIRoutes...)...)...,
)

// requiresAuthentication checks if a route requires authentication
func requiresAuthentication(path string) bool {
	// Check exact matches first
	for _, route := range ProtectedRoutes {
		if path == route {
			return true
		}
	}

	// Check admin routes (prefix matching)
	for _, route := range AdminRoutes {
		if len(path) >= len(route) && path[:len(route)] == route {
			return true // All admin API routes are protected
		}
	}

	return false
}

// isPublicRoute checks if a route is explicitly public
func isPublicRoute(path string) bool {
	for _, route := range PublicRoutes {
		if path == route {
			return true
		}
	}

	// Check OAuth routes (prefix matching for /auth/*)
	for _, route := range OAuthRoutes {
		if len(path) >= len(route) && path[:len(route)] == route {
			return true
		}
	}

	return false
}

// isAuthAPIRoute checks if a route is an auth API that might need session info
func isAuthAPIRoute(path string) bool {
	for _, route := range AuthAPIRoutes {
		if path == route {
			return true
		}
	}
	return false
}

// GetRouteCategory returns the category of a route for debugging
func GetRouteCategory(path string) string {
	if requiresAuthentication(path) {
		return "PROTECTED"
	}
	if isPublicRoute(path) {
		return "PUBLIC"
	}
	if isAuthAPIRoute(path) {
		return "AUTH_API"
	}
	return "UNKNOWN"
}


// RouteInfo provides information about application routes
type RouteInfo struct {
	Name        string `json:"name"`
	Method      string `json:"method"`
	Pattern     string `json:"pattern"`
	Description string `json:"description"`
}

// GetAllRoutes returns information about all application routes
func GetAllRoutes() []RouteInfo {
	return []RouteInfo{
		// Public Routes
		{Name: "home", Method: "GET", Pattern: "/", Description: "Main landing page"},
		{Name: "health", Method: "GET", Pattern: "/health", Description: "Health check endpoint"},
		{Name: "login", Method: "GET", Pattern: "/login", Description: "Login page"},

		// OAuth Routes
		{Name: "oauth_login", Method: "GET", Pattern: "/auth/login", Description: "OAuth provider login"},
		{Name: "oauth_callback", Method: "GET", Pattern: "/auth/callback", Description: "OAuth callback handler"},

		// Protected Routes
		{Name: "profile", Method: "GET", Pattern: "/profile", Description: "User profile page"},

		// Admin Routes
		{Name: "admin_dashboard", Method: "GET", Pattern: "/admin", Description: "Admin dashboard"},
		{Name: "admin_get_users", Method: "GET", Pattern: "/api/admin/users", Description: "Get users API"},
		{Name: "admin_get_analytics", Method: "GET", Pattern: "/api/admin/analytics", Description: "Get analytics API"},
		{Name: "admin_get_settings", Method: "GET", Pattern: "/api/admin/settings", Description: "Get settings API"},
		{Name: "admin_get_logs", Method: "GET", Pattern: "/api/admin/logs", Description: "Get logs API"},

		// Auth API Routes
		{Name: "logout", Method: "POST", Pattern: "/api/auth/logout", Description: "User logout"},
		{Name: "set_session", Method: "POST", Pattern: "/api/auth/set-session", Description: "Set session"},
	}
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
		TotalRoutes:     13,
		PublicRoutes:    3,
		ProtectedRoutes: 1,
		AdminRoutes:     5,
		AuthAPIRoutes:   4,
	}
}
