package routes

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

// RequiresAuthentication checks if a route requires authentication
func RequiresAuthentication(path string) bool {
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

// IsPublicRoute checks if a route is explicitly public
func IsPublicRoute(path string) bool {
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

// IsAuthAPIRoute checks if a route is an auth API that might need session info
func IsAuthAPIRoute(path string) bool {
	for _, route := range AuthAPIRoutes {
		if path == route {
			return true
		}
	}
	return false
}

// GetRouteCategory returns the category of a route for debugging
func GetRouteCategory(path string) string {
	if RequiresAuthentication(path) {
		return "PROTECTED"
	}
	if IsPublicRoute(path) {
		return "PUBLIC"
	}
	if IsAuthAPIRoute(path) {
		return "AUTH_API"
	}
	return "UNKNOWN"
}