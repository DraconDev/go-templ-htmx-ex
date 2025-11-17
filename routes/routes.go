package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/DraconDev/go-templ-htmx-ex/handlers"
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
}

// SetupRoutes configures and returns the router with all routes
// This version doesn't use the RouteDefinitions to avoid compilation issues
// with handler instance methods
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

	// These routes are configured directly in main.go since they need handler instances

	// =============================================================================
	// PROTECTED USER ROUTES - Authentication required
	// =============================================================================

	// User profile page - Display user information and account details
	router.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")

	// =============================================================================
	// ADMIN ROUTES - Admin authentication required
	// =============================================================================

	// These routes are configured directly in main.go since they need handler instances

	// =============================================================================
	// SESSION MANAGEMENT API - Authentication required
	// =============================================================================

	// These routes are configured directly in main.go since they need handler instances

	// Static files (for CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return router
}