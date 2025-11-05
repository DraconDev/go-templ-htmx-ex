package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AuthClient handles communication with the authentication service
type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

// AuthResponse represents the standard response from auth service
type AuthResponse struct {
	UserID       string `json:"user_id"`
	SessionToken string `json:"session_token"`
	Email        string `json:"email,omitempty"`
	Valid        bool   `json:"valid,omitempty"`
	ProjectIDs   []string `json:"project_ids,omitempty"`
	Error        string `json:"error,omitempty"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	ProjectID  string `json:"project_id"`
}

// NewAuthClient creates a new authentication client
func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Login attempts to authenticate a user and returns JWT token
func (c *AuthClient) Login(email, password string) (*AuthResponse, error) {
	loginReq := LoginRequest{
		Email:    email,
		Password: password,
	}

	return c.makeRequest("/api/auth/login", "POST", loginReq)
}

// Register creates a new user account
func (c *AuthClient) Register(email, password, projectID string) (*AuthResponse, error) {
	registerReq := RegisterRequest{
		Email:     email,
		Password:  password,
		ProjectID: projectID,
	}

	return c.makeRequest("/api/auth/register", "POST", registerReq)
}

// ValidateSession checks if a session token is valid
func (c *AuthClient) ValidateSession(sessionToken string) (*AuthResponse, error) {
	type ValidateRequest struct {
		SessionToken string `json:"session_token"`
	}
	
	validateReq := ValidateRequest{
		SessionToken: sessionToken,
	}

	return c.makeRequest("/api/auth/validate", "POST", validateReq)
}

// HealthCheck checks if the auth service is available
func (c *AuthClient) HealthCheck() (*AuthResponse, error) {
	return c.makeRequest("/api/health", "GET", nil)
}

// makeRequest handles the HTTP communication with the auth service
func (c *AuthClient) makeRequest(endpoint, method string, data interface{}) (*AuthResponse, error) {
	var reqBody io.Reader
	var err error

	// Prepare request body
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = nil
	}

	// Create HTTP request
	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Go-Auth-Client/1.0")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode >= 400 {
		if authResp.Error == "" {
			authResp.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body))
		}
	}

	return &authResp, nil
}
