package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

)

func TestGetRouteCategory(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/", "PUBLIC"},
		{"/login", "PUBLIC"},
		{"/health", "PUBLIC"},
		{"/auth/callback", "PUBLIC"},
		{"/auth/login", "PUBLIC"},
		{"/test", "PUBLIC"},
		{"/api/auth/exchange-code", "AUTH_API"},
		{"/api/auth/set-session", "AUTH_API"},
		{"/api/auth/logout", "AUTH_API"},
		{"/profile", "PROTECTED"},
		{"/admin", "PROTECTED"},
		{"/api/admin/users", "PROTECTED"},
		{"/api/admin/analytics", "PROTECTED"},
		{"/unknown/route", "UNKNOWN"},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			result := getRouteCategory(test.path)
			if result != test.expected {
				t.Errorf("getRouteCategory(%q) = %q, want %q", test.path, result, test.expected)
			}
		})
	}
}

func TestRequiresAuthentication(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		// Public routes - should NOT require authentication
		{"/", false},
		{"/login", false},
		{"/health", false},
		{"/test", false},
		{"/auth/callback", false},
		{"/auth/login", false},
		
		// Auth API routes - should NOT require authentication
		{"/api/auth/exchange-code", false},
		{"/api/auth/set-session", false},
		{"/api/auth/logout", false},
		
		// Protected routes - SHOULD require authentication
		{"/profile", true},
		{"/admin", true},
		{"/api/admin/users", true},
		{"/api/admin/analytics", true},
		{"/api/admin/settings", true},
		{"/api/admin/logs", true},
		
		// Unknown routes - should NOT require authentication
		{"/unknown", false},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			result := requiresAuthentication(test.path)
			if result != test.expected {
				t.Errorf("requiresAuthentication(%q) = %v, want %v", test.path, result, test.expected)
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		s       string
		prefix  string
		expected bool
	}{
		{"hello world", "hello", true},
		{"hello world", "world", false},
		{"hello world", "hello world", true},
		{"hello world", "hello world!", false},
		{"", "prefix", false},
		{"prefix", "", true},
		{"", "", true},
	}

	for _, test := range tests {
		t.Run(test.s+"_"+test.prefix, func(t *testing.T) {
			result := hasPrefix(test.s, test.prefix)
			if result != test.expected {
				t.Errorf("hasPrefix(%q, %q) = %v, want %v", test.s, test.prefix, result, test.expected)
			}
		})
	}
}

func TestAuthMiddlewarePublicRoutes(t *testing.T) {
	// Test that public routes don't redirect
	publicRoutes := []string{"/", "/login", "/health", "/auth/callback", "/test"}

	for _, route := range publicRoutes {
		t.Run("public_"+route, func(t *testing.T) {
			req := httptest.NewRequest("GET", route, nil)
			rr := httptest.NewRecorder()

			// Create a simple test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			// Apply middleware
			middleware := AuthMiddleware(testHandler)
			middleware.ServeHTTP(rr, req)

			// Should not redirect for public routes
			if rr.Code == http.StatusFound {
				t.Errorf("Public route %q should not redirect, got status %d", route, rr.Code)
			}
		})
	}
}

func TestAuthMiddlewareAuthAPIRoutes(t *testing.T) {
	// Test that auth API routes don't require authentication
	authAPIRoutes := []string{"/api/auth/exchange-code", "/api/auth/set-session", "/api/auth/logout"}

	for _, route := range authAPIRoutes {
		t.Run("auth_api_"+route, func(t *testing.T) {
			req := httptest.NewRequest("POST", route, nil)
			rr := httptest.NewRecorder()

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			middleware := AuthMiddleware(testHandler)
			middleware.ServeHTTP(rr, req)

			// Should not redirect and should not return 401
			if rr.Code == http.StatusFound {
				t.Errorf("Auth API route %q should not redirect, got status %d", route, rr.Code)
			}
			if rr.Code == http.StatusUnauthorized {
				t.Errorf("Auth API route %q should not require authentication, got status %d", route, rr.Code)
			}
		})
	}
}

func TestAuthMiddlewareProtectedRoutes(t *testing.T) {
	// Test that protected routes redirect when no session
	protectedRoutes := []string{"/profile", "/admin", "/api/admin/users"}

	for _, route := range protectedRoutes {
		t.Run("protected_"+route, func(t *testing.T) {
			req := httptest.NewRequest("GET", route, nil)
			rr := httptest.NewRecorder()

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			middleware := AuthMiddleware(testHandler)
			middleware.ServeHTTP(rr, req)

			// Should redirect to login for web routes
			if route == "/profile" || route == "/admin" {
				if rr.Code != http.StatusFound {
					t.Errorf("Protected web route %q should redirect to login, got status %d", route, rr.Code)
				}
				if location := rr.Header().Get("Location"); location != "/login" {
					t.Errorf("Expected redirect to /login, got %q", location)
				}
			}

			// Should return 401 for API routes
			if route == "/api/admin/users" {
				if rr.Code != http.StatusUnauthorized {
					t.Errorf("Protected API route %q should return 401, got status %d", route, rr.Code)
				}
				if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
					t.Errorf("Expected JSON response, got %q", contentType)
				}
			}
		})
	}
}

func TestGetUserFromContext(t *testing.T) {
	// Test getting user info from request context
	req := httptest.NewRequest("GET", "/", nil)

	// Test with no user context
	userInfo := GetUserFromContext(req)
	if userInfo.LoggedIn {
		t.Error("Expected user to not be logged in with no context")
	}

	// Test with user context (context setup would be done by middleware)
	req = req.WithContext(req.Context())

	// Note: In actual usage, the user context is set by the middleware
	// This test shows the expected behavior
	userInfo = GetUserFromContext(req)
	if userInfo.LoggedIn {
		t.Error("Expected user to not be logged in without proper context setup")
	}
}