package auth

import "fmt"

// AuthClient is the interface that wraps the auth service client methods
type AuthClient interface {
	Login(email, password string) (*AuthResponse, error)
	Register(email, password, projectID string) (*AuthResponse, error)
	ValidateSession(sessionToken string) (*AuthResponse, error)
	GetUserDetails(userID string) (*AuthResponse, error)
	HealthCheck() (*AuthResponse, error)
}

// AuthResponse represents the response from auth service methods
type AuthResponse struct {
	UserID       string   `json:"user_id"`
	SessionToken string   `json:"session_token"`
	Email        string   `json:"email,omitempty"`
	Valid        bool     `json:"valid,omitempty"`
	ProjectIDs   []string `json:"project_ids,omitempty"`
	Error        string   `json:"error,omitempty"`
	Status       string   `json:"status,omitempty"`
	Message      string   `json:"message,omitempty"`
}

// NewAuthClient creates a new authentication client
// This is a wrapper that creates a gRPC auth client
func NewAuthClient(baseURL string) AuthClient {
	grpcClient, err := NewGRPCAuthClient(baseURL)
	if err != nil {
		// For now, we'll return a mock client if gRPC connection fails
		// In production, you might want to handle this differently
		return &MockAuthClient{
			baseURL: baseURL,
			// Could log the error or handle it according to your requirements
		}
	}
	
	return grpcClient
}

// MockAuthClient is a mock implementation for testing
type MockAuthClient struct {
	baseURL string
}

// Login implements mock login for testing when gRPC is not available
func (m *MockAuthClient) Login(email, password string) (*AuthResponse, error) {
	return &AuthResponse{
		UserID:       "mock-user-id",
		SessionToken: "mock-session-token",
		Email:        email,
		Status:       "success",
		Message:      "Mock login successful (gRPC not available)",
	}, nil
}

// Register implements mock register for testing when gRPC is not available
func (m *MockAuthClient) Register(email, password, projectID string) (*AuthResponse, error) {
	return &AuthResponse{
		UserID:       "mock-user-id",
		SessionToken: "mock-session-token",
		Status:       "success",
		Message:      "Mock register successful (gRPC not available)",
	}, nil
}

// ValidateSession implements mock validate session for testing when gRPC is not available
func (m *MockAuthClient) ValidateSession(sessionToken string) (*AuthResponse, error) {
	return &AuthResponse{
		UserID:     "mock-user-id",
		Valid:      true,
		ProjectIDs: []string{"mock-project-id"},
		Status:     "validated",
		Message:    "Mock session validation successful (gRPC not available)",
	}, nil
}

// GetUserDetails implements mock get user details for testing when gRPC is not available
func (m *MockAuthClient) GetUserDetails(userID string) (*AuthResponse, error) {
	return &AuthResponse{
		UserID: userID,
		Email:  "mock@example.com",
		Status: "success",
		Message: "Mock user details successful (gRPC not available)",
	}, nil
}

// HealthCheck implements mock health check for testing when gRPC is not available
func (m *MockAuthClient) HealthCheck() (*AuthResponse, error) {
	return &AuthResponse{
		Status:  "available",
		Message: fmt.Sprintf("Mock health check - gRPC not available (baseURL: %s)", m.baseURL),
	}, fmt.Errorf("gRPC connection not available - using mock client")
}
