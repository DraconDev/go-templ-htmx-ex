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

// RouteConfig defines all application routes with proper organization
type RouteConfig struct {
	PublicRoutes      []Route
	ProtectedRoutes   []Route
	AdminRoutes       []Route
	AuthAPIRoutes     []Route
}

// GetRoutes returns all application routes
func GetRoutes() RouteConfig {
	return RouteConfig{
		// =============================================================================
		// PUBLIC ROUTES - No authentication required
		// =============================================================================
		PublicRoutes: []Route{
			{
				Name:        "home",
				Method:      "GET",
				Pattern:     "/",
				HandlerFunc: handlers.HomeHandler,
			},
			{
				Name:        "health",
				Method:      "GET",
				Pattern:     "/health",
				HandlerFunc: handlers.HealthHandler,
			},
			{
				Name:        "login",
				Method:      "GET",
				Pattern:     "/login",
				HandlerFunc: handlers.LoginHandler,
			},
		},

		// =============================================================================
		// PROTECTED USER ROUTES - Authentication required
		// =============================================================================
		ProtectedRoutes: []Route{
			{
				Name:        "profile",
				Method:      "GET",
				Pattern:     "/profile",
				HandlerFunc: handlers.ProfileHandler,
			},
		},

		// =============================================================================
		// ADMIN ROUTES - Admin authentication required
		// =============================================================================
		AdminRoutes: []Route{
			{
				Name:        "admin_dashboard",
				Method:      "GET",
				Pattern:     "/admin",
				HandlerFunc: admin.AdminDashboardHandler,
			},
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
		},

		// =============================================================================
		// AUTH API ROUTES - Handle authentication endpoints
		// =============================================================================
		AuthAPIRoutes: []Route{
			{
				Name:        "oauth_login",
				Method:      "GET",
				Pattern:     "/auth/login",
				HandlerFunc: authHandlers.AuthHandler.LoginHandler,
			},
			{
				Name:        "oauth_callback",
				Method:      "GET",
				Pattern:     "/auth/callback",
				HandlerFunc: authHandlers.AuthHandler.AuthCallbackHandler,
			},
			{
				Name:        "logout",
				Method:      "POST",
				Pattern:     "/api/auth/logout",
				HandlerFunc: authHandlers.AuthHandler.LogoutHandler,
			},
			{
				Name:        "set_session",
				Method:      "POST",
				Pattern:     "/api/auth/set-session",
				HandlerFunc: authHandlers.AuthHandler.SetSessionHandler,
			},
		},
	}
}

// SetupRoutes configures and returns the router with all routes
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Add authentication middleware to all routes
	router.Use(middleware.AuthMiddleware)

	// Get all route definitions
	routes := GetRoutes()

	// Register all public routes
	for _, route := range routes.PublicRoutes {
		router.HandleFunc(route.Pattern, route.HandlerFunc).
			Methods(route.Method).
			Name(route.Name)
	}

	// Register all protected routes
	for _, route := range routes.ProtectedRoutes {
		router.HandleFunc(route.Pattern, route.HandlerFunc).
			Methods(route.Method).
			Name(route.Name)
	}

	// Register all admin routes
	for _, route := range routes.AdminRoutes {
		router.HandleFunc(route.Pattern, route.HandlerFunc).
			Methods(route.Method).
			Name(route.Name)
	}

	// Register all auth API routes
	for _, route := range routes.AuthAPIRoutes {
		router.HandleFunc(route.Pattern, route.HandlerFunc).
			Methods(route.Method).
			Name(route.Name)
	}

	// Static files (for CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return router
}

// GetPublicRoutes returns just the public routes for external use
func GetPublicRoutes() []Route {
	return GetRoutes().PublicRoutes
}

// GetProtectedRoutes returns just the protected routes for external use
func GetProtectedRoutes() []Route {
	return GetRoutes().ProtectedRoutes
}

// GetAdminRoutes returns just the admin routes for external use
func GetAdminRoutes() []Route {
	return GetRoutes().AdminRoutes
}

// GetAuthAPIRoutes returns just the auth API routes for external use
func GetAuthAPIRoutes() []Route {
	return GetRoutes().AuthAPIRoutes
}