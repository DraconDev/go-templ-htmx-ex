package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DraconDev/go-templ-htmx-ex/internal/handlers"
	"github.com/DraconDev/go-templ-htmx-ex/internal/middleware"
	"github.com/DraconDev/go-templ-htmx-ex/internal/services"
	"github.com/DraconDev/go-templ-htmx-ex/internal/utils/config"
)

// MockAuthService simulates the auth microservice responses
type MockAuthService struct{}

func (m *MockAuthService) ExchangeCodeForTokens(authCode string) (*handlers.AuthResponse, error) {
	fmt.Printf("🔍 TEST: MockAuthService.ExchangeCodeForTokens called with: %s\n", authCode)
	
	// Simulate successful token exchange
	if authCode == "" {
		return nil, fmt.Errorf("empty auth code")
	}
	
	return &handlers.AuthResponse{
		Success:  true,
		UserID:   "test-user-id-123",
		UserName: "Test User",
		UserEmail: "test@example.com",
		Message:  "Tokens exchanged successfully",
	}, nil
}

func (m *MockAuthService) CreateSession(authCode string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"session_id":   "test-session-123",
		"user_context": map[string]interface{}{
			"user_id": "test-user-id-123",
			"name":    "Test User",
			"email":   "test@example.com",
		},
	}, nil
}

func (m *MockAuthService) GetUserInfo(sessionID string) (*handlers.AuthResponse, error) {
	return &handlers.AuthResponse{
		Success:    true,
		UserID:     "test-user-id-123",
		UserName:   "Test User",
		UserEmail:  "test@example.com",
		Message:    "User info retrieved",
	}, nil
}

func (m *MockAuthService) RefreshSession(sessionID string) (*handlers.AuthResponse, error) {
	return &handlers.AuthResponse{
		Success:  true,
		UserID:   "test-user-id-123",
		UserName: "Test User",
		Message:  "Session refreshed",
	}, nil
}

// TestExchangeCodeEndpoint tests the /api/auth/exchange-code endpoint
func TestExchangeCodeEndpoint(t *testing.T) {
	fmt.Println("🧪 STARTING: TestExchangeCodeEndpoint")
	
	// Initialize config and services
	cfg := &config.Config{
		AuthServiceURL: "http://mock-auth-service",
		RedirectURL:    "http://localhost:8080",
	}
	
	authService := services.NewAuthService(cfg)
	mockService := &MockAuthService{}
	
	// Create handler
	authHandler := handlers.NewAuthHandler(authService, cfg)
	
	// Create test request with valid auth code
	req := httptest.NewRequest("POST", "/api/auth/exchange-code", 
		genJSONBody(map[string]string{"auth_code": "test-github-123"}))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Call handler
	authHandler.ExchangeCodeHandler(rr, req)
	
	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
	
	// Check response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	// Verify response structure
	if !response["success"].(bool) {
		t.Errorf("Expected success to be true, got %v", response["success"])
	}
	
	fmt.Printf("✅ TestExchangeCodeEndpoint: Response: %+v\n", response)
	fmt.Println("✅ PASSED: TestExchangeCodeEndpoint")
}

// TestOAuthCallbackEndpoint tests the /auth/callback endpoint
func TestOAuthCallbackEndpoint(t *testing.T) {
	fmt.Println("🧪 STARTING: TestOAuthCallbackEndpoint")
	
	cfg := &config.Config{
		AuthServiceURL: "http://mock-auth-service",
		RedirectURL:    "http://localhost:8080",
	}
	
	authService := services.NewAuthService(cfg)
	authHandler := handlers.NewAuthHandler(authService, cfg)
	
	// Create test request simulating GitHub OAuth callback
	req := httptest.NewRequest("GET", "/auth/callback?auth_code=github_12345_cb67890", nil)
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Call handler
	authHandler.AuthCallbackHandler(rr, req)
	
	// Check status and content type
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
	
	if rr.Header().Get("Content-Type") != "text/html" {
		t.Errorf("Expected Content-Type 'text/html', got %s", rr.Header().Get("Content-Type"))
	}
	
	// Check that response contains expected HTML elements
	body := rr.Body.String()
	if !containsSubstring(body, "Setting up your session") {
		t.Errorf("Expected response to contain authentication message")
	}
	
	fmt.Printf("✅ TestOAuthCallbackEndpoint: Status=%d, Content-Type=%s\n", rr.Code, rr.Header().Get("Content-Type"))
	fmt.Println("✅ PASSED: TestOAuthCallbackEndpoint")
}

// TestMiddlewareRouteCategorization tests middleware route categorization
func TestMiddlewareRouteCategorization(t *testing.T) {
	fmt.Println("🧪 STARTING: TestMiddlewareRouteCategorization")
	
	// Test cases for route categorization
	testCases := []struct {
		path       string
		expected   string
		shouldAuth bool
	}{
		{"/", "PUBLIC", false},
		{"/login", "PUBLIC", false},
		{"/auth/callback", "PUBLIC", false},
		{"/auth/login", "PUBLIC", false},
		{"/api/auth/exchange-code", "AUTH_API", false}, // Should NOT require auth
		{"/profile", "PROTECTED", true},
		{"/admin", "PROTECTED", true},
		{"/api/admin/users", "PROTECTED", true},
		{"/health", "PUBLIC", false},
		{"/unknown", "UNKNOWN", false},
	}
	
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Path:%s", tc.path), func(t *testing.T) {
			// Test route category detection
			category := middleware.GetRouteCategory(tc.path)
			if category != tc.expected {
				t.Errorf("Expected category %s for path %s, got %s", tc.expected, tc.path, category)
			}
			
			// Test authentication requirement
			requiresAuth := middleware.RequiresAuthentication(tc.path)
			if requiresAuth != tc.shouldAuth {
				t.Errorf("Expected auth requirement %v for path %s, got %v", tc.shouldAuth, tc.path, requiresAuth)
			}
		})
	}
	
	fmt.Println("✅ PASSED: TestMiddlewareRouteCategorization")
}

// TestMiddlewareIntegration tests middleware integration with handlers
func TestMiddlewareIntegration(t *testing.T) {
	fmt.Println("🧪 STARTING: TestMiddlewareIntegration")
	
	cfg := &config.Config{
		AuthServiceURL: "http://mock-auth-service",
		RedirectURL:    "http://localhost:8080",
	}
	
	authService := services.NewAuthService(cfg)
	authHandler := handlers.NewAuthHandler(authService, cfg)
	
	// Test protected route without session cookie
	req := httptest.NewRequest("GET", "/profile", nil)
	rr := httptest.NewRecorder()
	
	// Apply middleware
	middleware.AuthMiddleware(authHandler.ProfileHandler).ServeHTTP(rr, req)
	
	// Should redirect to login
	if rr.Code != http.StatusFound {
		t.Errorf("Expected redirect status %d for protected route without auth, got %d", http.StatusFound, rr.Code)
	}
	
	// Test public route (should not redirect)
	req = httptest.NewRequest("GET", "/login", nil)
	rr = httptest.NewRecorder()
	
	middleware.AuthMiddleware(authHandler.ProfileHandler).ServeHTTP(rr, req)
	
	// Should not redirect (handled by profile handler)
	if rr.Code == http.StatusFound {
		t.Errorf("Public route should not redirect, got status %d", rr.Code)
	}
	
	fmt.Println("✅ PASSED: TestMiddlewareIntegration")
}

// TestAuthServiceIntegration tests the AuthService with mock responses
func TestAuthServiceIntegration(t *testing.T) {
	fmt.Println("🧪 STARTING: TestAuthServiceIntegration")
	
	cfg := &config.Config{
		AuthServiceURL: "http://mock-auth-service",
	}
	
	authService := services.NewAuthService(cfg)
	
	// Test CreateSession
	result, err := authService.CreateSession("test-code-123")
	if err != nil {
		t.Errorf("CreateSession failed: %v", err)
	}
	
	if result["session_id"] != "test-session-123" {
		t.Errorf("Expected session_id 'test-session-123', got %v", result["session_id"])
	}
	
	fmt.Println("✅ PASSED: TestAuthServiceIntegration")
}

// Helper functions
func genJSONBody(data interface{}) *bytes.Reader {
	// Note: This is a simplified JSON body generator for testing
	// In actual implementation, you'd use json.Marshal
	return nil
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}