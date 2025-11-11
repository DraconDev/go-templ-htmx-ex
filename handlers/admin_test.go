package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
)

// MockQueries provides minimal mock implementation for testing
type MockQueries struct {
	TotalUsers      int64
	SignupsToday    int64
	SignupsThisWeek int64
	RecentUsers     []mockUser
}

type mockUser struct {
	ID       int64
	Email    string
	Name     string
	CreatedAt time.Time
}

func (m *MockQueries) CountUsers(ctx context.Context) (int64, error) {
	return m.TotalUsers, nil
}

func (m *MockQueries) CountUsersCreatedToday(ctx context.Context) (int64, error) {
	return m.SignupsToday, nil
}

func (m *MockQueries) CountUsersCreatedThisWeek(ctx context.Context) (int64, error) {
	return m.SignupsThisWeek, nil
}

func (m *MockQueries) GetRecentUsers(ctx context.Context) ([]struct {
	ID       int64
	Email    string
	Name     string
	CreatedAt time.Time
}, error) {
	result := make([]struct {
		ID       int64
		Email    string
		Name     string
		CreatedAt time.Time
	}, len(m.RecentUsers))
	
	for i, user := range m.RecentUsers {
		result[i] = struct {
			ID       int64
			Email    string
			Name     string
			CreatedAt time.Time
		}{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			CreatedAt: user.CreatedAt,
		}
	}
	return result, nil
}

// CreateAdminHandlerForTest creates an admin handler with mock data for testing
func CreateAdminHandlerForTest() *AdminHandler {
	mockQueries := &MockQueries{
		TotalUsers:      42,
		SignupsToday:    5,
		SignupsThisWeek: 12,
		RecentUsers: []mockUser{
			{
				ID:       1,
				Email:    "test@example.com",
				Name:     "Test User",
				CreatedAt: time.Now(),
			},
			{
				ID:       2,
				Email:    "admin@example.com",
				Name:     "Admin User",
				CreatedAt: time.Now(),
			},
		},
	}

	return &AdminHandler{
		Config:  &config.Config{},
		Queries: mockQueries,
	}
}

// TestAdminDashboardAccess tests admin dashboard access with valid admin
func TestAdminDashboardAccess(t *testing.T) {
	handler := CreateAdminHandlerForTest()
	
	// Create test request
	req := httptest.NewRequest("GET", "/admin", nil)
	
	// Set admin user claims in context (simulating middleware)
	ctx := req.Context()
	ctx = setTestUserClaims(ctx, "admin@example.com", "Admin User", true)
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Execute handler
	handler.AdminDashboardHandler(rr, req.WithContext(ctx))
	
	// Verify response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
	
	body := rr.Body.String()
	
	// Check that dashboard contains expected content
	if !strings.Contains(body, "Admin Dashboard") {
		t.Errorf("Expected dashboard title not found in response")
	}
	
	if !strings.Contains(body, "42") {
		t.Errorf("Expected total users count not found")
	}
	
	if !strings.Contains(body, "5") {
		t.Errorf("Expected today's signups not found")
	}
	
	if !strings.Contains(body, "12") {
		t.Errorf("Expected this week's signups not found")
	}
}

// TestAdminDashboardUnauthorized tests admin dashboard without admin privileges
func TestAdminDashboardUnauthorized(t *testing.T) {
	handler := CreateAdminHandlerForTest()
	
	// Create test request with non-admin user
	req := httptest.NewRequest("GET", "/admin", nil)
	
	// Set non-admin user claims
	ctx := req.Context()
	ctx = setTestUserClaims(ctx, "user@example.com", "Regular User", false)
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Execute handler
	handler.AdminDashboardHandler(rr, req.WithContext(ctx))
	
	// Should redirect for non-admin users
	if rr.Code != http.StatusFound {
		t.Errorf("Expected redirect status %d, got %d", http.StatusFound, rr.Code)
	}
	
	location := rr.Header().Get("Location")
	if location != "/" {
		t.Errorf("Expected redirect to '/', got '%s'", location)
	}
}

// TestAdminDashboardNoAuth tests admin dashboard without authentication
func TestAdminDashboardNoAuth(t *testing.T) {
	handler := CreateAdminHandlerForTest()
	
	// Create test request without authentication
	req := httptest.NewRequest("GET", "/admin", nil)
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Execute handler
	handler.AdminDashboardHandler(rr, req)
	
	// Should redirect for unauthenticated users
	if rr.Code != http.StatusFound {
		t.Errorf("Expected redirect status %d, got %d", http.StatusFound, rr.Code)
	}
	
	location := rr.Header().Get("Location")
	if location != "/" {
		t.Errorf("Expected redirect to '/', got '%s'", location)
	}
}

// TestAdminDashboardDataLoading tests that dashboard loads mock data correctly
func TestAdminDashboardDataLoading(t *testing.T) {
	handler := CreateAdminHandlerForTest()
	
	req := httptest.NewRequest("GET", "/admin", nil)
	ctx := req.Context()
	ctx = setTestUserClaims(ctx, "admin@example.com", "Admin User", true)
	
	rr := httptest.NewRecorder()
	handler.AdminDashboardHandler(rr, req.WithContext(ctx))
	
	// Verify dashboard contains expected mock data
	body := rr.Body.String()
	
	// Check specific data points
	if !strings.Contains(body, "42") {
		t.Error("Expected total users count '42' not found in dashboard")
	}
	
	if !strings.Contains(body, "5") {
		t.Error("Expected today's signups '5' not found in dashboard")
	}
	
	if !strings.Contains(body, "12") {
		t.Error("Expected this week's signups '12' not found in dashboard")
	}
	
	if !strings.Contains(body, "Test User") {
		t.Error("Expected recent user 'Test User' not found in dashboard")
	}
	
	if !strings.Contains(body, "Admin User") {
		t.Error("Expected recent user 'Admin User' not found in dashboard")
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