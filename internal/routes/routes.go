package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/DraconDev/go-templ-htmx-ex/internal/routing"
)

// HandlerInstances holds all handler instances for route registration
type HandlerInstances struct {
	// These will be populated with actual handler instances
}

// SetupRoutes configures and returns the router with all routes
// This approach accepts handler instances to avoid method reference issues
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Static files (for CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return router
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
