package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/db/sqlc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAdminHandler creates a test admin handler with mock database
func MockAdminHandler() *AdminHandler {
	return &AdminHandler{
		Queries: &mockQueries{},
	}
}

// mockQueries implements the sqlc.Queries interface for testing
type mockQueries struct{}

func (m *mockQueries) CountUsers(ctx context.Context) (int64, error) {
	return 42, nil
}

func (m *mockQueries) CountUsersCreatedToday(ctx context.Context) (int64, error) {
	return 5, nil
}

func (m *mockQueries) CountUsersCreatedThisWeek(ctx context.Context) (int64, error) {
	return 12, nil
}

func (m *mockQueries) GetRecentUsers(ctx context.Context) ([]sqlc.GetRecentUsersRow, error) {
	return []sqlc.GetRecentUsersRow{
		{
			ID:    1,
			Email: "test@example.com",
			Name:  "Test User",
			CreatedAt: time.Now(),
		},
		{
			ID:    2,
			Email: "admin@example.com", 
			Name:  "Admin User",
			CreatedAt: time.Now(),
		},
	}, nil
}

// TestAdminDashboardAccess tests admin dashboard access
func TestAdminDashboardAccess(t *testing.T) {
	handler := MockAdminHandler()
	
	// Create test request
	req := httptest.NewRequest("GET", "/admin", nil)
	
	// Set JWT claims in context (simulating middleware)
	ctx := req.Context()
	ctx = setUserClaims(ctx, "test@example.com", "Test User", true)
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Execute handler
	handler.AdminDashboardHandler(rr, req.WithContext(ctx))
	
	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Admin Dashboard")
	assert.Contains(t, rr.Body.String(), "42") // Total users from mock
}

// TestAdminDashboardUnauthorized tests admin dashboard without admin privileges
func TestAdminDashboardUnauthorized(t *testing.T) {
	handler := MockAdminHandler()
	
	// Create test request with non-admin user
	req := httptest.NewRequest("GET", "/admin", nil)
	
	// Set non-admin user claims
	ctx := req.Context()
	ctx = setUserClaims(ctx, "regular@example.com", "Regular User", false)
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Execute handler
	handler.AdminDashboardHandler(rr, req.WithContext(ctx))
	
	// Should redirect for non-admin users
	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Contains(t, rr.Header().Get("Location"), "/")
}

// TestAdminDashboardNoAuth tests admin dashboard without authentication
func TestAdminDashboardNoAuth(t *testing.T) {
	handler := MockAdminHandler()
	
	// Create test request without authentication
	req := httptest.NewRequest("GET, "/admin", nil)
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Execute handler
	handler.AdminDashboardHandler(rr, req)
	
	// Should redirect for unauthenticated users
	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Contains(t, rr.Header().Get("Location"), "/")
}

// setUserClaims is a helper to set JWT claims in request context for testing
func setUserClaims(ctx context.Context, email, name string, isAdmin bool) context.Context {
	claims := map[string]interface{}{
		"email":    email,
		"name":     name,
		"is_admin": isAdmin,
		"sub":      "test-user-id",
		"iss":      "auth-ms",
	}
	return context.WithValue(ctx, "user_claims", claims)
}

// TestAdminDashboardDataLoading tests that dashboard loads real data
func TestAdminDashboardDataLoading(t *testing.T) {
	handler := MockAdminHandler()
	
	req := httptest.NewRequest("GET", "/admin", nil)
	ctx := req.Context()
	ctx = setUserClaims(ctx, "admin@example.com", "Admin User", true)
	
	rr := httptest.NewRecorder()
	handler.AdminDashboardHandler(rr, req.WithContext(ctx))
	
	// Verify dashboard contains expected data
	body := rr.Body.String()
	assert.Contains(t, body, "42")   // Total users
	assert.Contains(t, body, "5")    // Today's signups
	assert.Contains(t, body, "12")   // This week's signups
	assert.Contains(t, body, "Test User")
	assert.Contains(t, body, "Admin User")
}