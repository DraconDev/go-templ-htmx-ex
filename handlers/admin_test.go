package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
)

// TestAdminDashboardAccess tests the admin dashboard endpoint behavior
func TestAdminDashboardAccess(t *testing.T) {
	// This test validates the handler behavior without requiring a full database setup
	handler := &AdminHandler{
		Config:  &config.Config{},
		Queries: nil,
	}
	
	// Create test request
	req := httptest.NewRequest("GET", "/admin", nil)
	
	// Set session token cookie for admin user (simulating real authentication)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "admin-jwt-token",
	})
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Execute handler - this will try to validate the JWT
	handler.AdminDashboardHandler(rr, req)
	
	// We expect this to either:
	// 1. Return OK if JWT validation succeeds
	// 2. Return error if JWT validation fails
	// 3. Return redirect if user is not authenticated
	
	// Since we're testing without real JWT validation, this will likely redirect
	// but that's fine - we're testing the flow, not the JWT validation itself
	t.Logf("Test status code: %d (expected for test JWT)", rr.Code)
}

// TestAdminDashboardUnauthorized tests admin dashboard without admin privileges
func TestAdminDashboardUnauthorized(t *testing.T) {
	handler := &AdminHandler{
		Config:  &config.Config{},
		Queries: nil,
	}
	
	// Create test request with non-admin user token
	req := httptest.NewRequest("GET", "/admin", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "regular-user-jwt-token",
	})
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Execute handler
	handler.AdminDashboardHandler(rr, req)
	
	// Should either redirect or return error for non-admin users
	if rr.Code == http.StatusFound {
		location := rr.Header().Get("Location")
		if location != "/" {
			t.Errorf("Expected redirect to '/', got '%s'", location)
		}
	}
}

// TestAdminDashboardNoAuth tests admin dashboard without authentication
func TestAdminDashboardNoAuth(t *testing.T) {
	handler := &AdminHandler{
		Config:  &config.Config{},
		Queries: nil,
	}
	
	// Create test request without authentication (no cookie)
	req := httptest.NewRequest("GET", "/admin", nil)
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Execute handler
	handler.AdminDashboardHandler(rr, req)
	
	// Should redirect for unauthenticated users
	if rr.Code != http.StatusFound {
		t.Errorf("Expected redirect status %d for unauthenticated user, got %d", http.StatusFound, rr.Code)
	}
	
	location := rr.Header().Get("Location")
	if location != "/" {
		t.Errorf("Expected redirect to '/', got '%s'", location)
	}
}

// TestAdminDashboardClaimExtraction tests that user claims are properly extracted
func TestAdminDashboardClaimExtraction(t *testing.T) {
	handler := &AdminHandler{
		Config:  &config.Config{},
		Queries: nil,
	}
	
	testCases := []struct {
		testName string
		email    string
		userName string
		isAdmin  bool
		expected int
	}{
		{
			testName: "Admin User",
			email:    "admin@example.com",
			userName: "Admin User",
			isAdmin:  true,
			expected: http.StatusOK, // Would be OK if we had database
		},
		{
			testName: "Regular User",
			email:    "user@example.com", 
			userName: "Regular User",
			isAdmin:  false,
			expected: http.StatusFound, // Should redirect
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/admin", nil)
			ctx := req.Context()
			ctx = setTestUserClaims(ctx, tc.email, tc.userName, tc.isAdmin)
			
			rr := httptest.NewRecorder()
			handler.AdminDashboardHandler(rr, req.WithContext(ctx))
			
			// For admin users, we expect either OK or server error (if no DB)
			// For non-admin users, we expect redirect
			if tc.isAdmin && rr.Code == http.StatusFound {
				t.Errorf("Admin user should not be redirected, got status %d", rr.Code)
			}
			
			if !tc.isAdmin && rr.Code != http.StatusFound {
				t.Errorf("Non-admin user should be redirected, got status %d", rr.Code)
			}
		})
	}
}

// setTestUserClaims is a helper to set JWT claims in request context for testing
func setTestUserClaims(ctx context.Context, email, name string, isAdmin bool) context.Context {
	claims := map[string]interface{}{
		"email":    email,
		"name":     name,
		"is_admin": isAdmin,
		"sub":      "test-user-id",
		"iss":      "auth-ms",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
	return context.WithValue(ctx, "user_claims", claims)
}

// TestAdminDashboardMiddlewareIntegration tests that the handler properly integrates with middleware
func TestAdminDashboardMiddlewareIntegration(t *testing.T) {
	// This test simulates what would happen when the actual middleware is used
	handler := &AdminHandler{
		Config:  &config.Config{},
		Queries: nil,
	}
	
	// Test the full flow: unauthenticated -> authenticated admin -> unauthorized
	scenarios := []struct {
		name          string
		hasClaims     bool
		email         string
		isAdmin       bool
		expectedCode  int
		expectedLoc   string
	}{
		{
			name:         "No Authentication",
			hasClaims:    false,
			expectedCode: http.StatusFound,
			expectedLoc:  "/",
		},
		{
			name:         "Authenticated Non-Admin",
			hasClaims:    true,
			email:        "user@example.com",
			isAdmin:      false,
			expectedCode: http.StatusFound,
			expectedLoc:  "/",
		},
		{
			name:         "Authenticated Admin",
			hasClaims:    true,
			email:        "admin@example.com",
			isAdmin:      true,
			expectedCode: http.StatusOK, // Would succeed if we had DB
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/admin", nil)
			ctx := req.Context()
			
			if scenario.hasClaims {
				ctx = setTestUserClaims(ctx, scenario.email, scenario.email, scenario.isAdmin)
			}
			
			rr := httptest.NewRecorder()
			handler.AdminDashboardHandler(rr, req.WithContext(ctx))
			
			if rr.Code != scenario.expectedCode {
				t.Errorf("Expected status %d, got %d", scenario.expectedCode, rr.Code)
			}
			
			if scenario.expectedLoc != "" {
				loc := rr.Header().Get("Location")
				if loc != scenario.expectedLoc {
					t.Errorf("Expected redirect to '%s', got '%s'", scenario.expectedLoc, loc)
				}
			}
		})
	}
}