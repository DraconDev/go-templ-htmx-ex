package auth

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/DraconDev/cerberus-auth-ms/proto"
)

// AuthClient handles communication with the authentication service via gRPC
type AuthClient struct {
	conn         *grpc.ClientConn
	authService  proto.AuthServiceClient
}

// AuthResponse represents the standard response from auth service
type AuthResponse struct {
	UserID       string   `json:"user_id"`
	SessionToken string   `json:"session_token"`
	Email        string   `json:"email,omitempty"`
	Valid        bool     `json:"valid,omitempty"`
	ProjectIDs   []string `json:"project_ids,omitempty"`
	Error        string   `json:"error,omitempty"`
}

// NewAuthClient creates a new authentication client with gRPC connection
func NewAuthClient(baseURL string) *AuthClient {
	// Create gRPC connection
	conn, err := grpc.Dial(baseURL, grpc.WithTransportCredentials(insecureCredentials()))
	if err != nil {
		fmt.Printf("Failed to connect to auth service: %v\n", err)
		return nil
	}

	// Create gRPC client
	authService := proto.NewAuthServiceClient(conn)

	return &AuthClient{
		conn:        conn,
		authService: authService,
	}
}

// insecureCredentials returns credentials that skip TLS verification (for development)
func insecureCredentials() grpc.DialOption {
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

// Login attempts to authenticate a user and returns JWT token
func (c *AuthClient) Login(email, password string) (*AuthResponse, error) {
	if c.authService == nil {
		return nil, fmt.Errorf("auth service not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create gRPC request
	loginReq := &proto.LoginRequest{
		Email:    email,
		Password: password,
	}

	// Make gRPC call
	loginResp, err := c.authService.Login(ctx, loginReq)
	if err != nil {
		return nil, fmt.Errorf("gRPC login failed: %w", err)
	}

	// Convert gRPC response to our response format
	return &AuthResponse{
		UserID:       loginResp.GetUserId(),
		SessionToken: loginResp.GetSessionToken(),
		Email:        email,
		Valid:        true,
	}, nil
}

// Register creates a new user account
func (c *AuthClient) Register(email, password, projectID string) (*AuthResponse, error) {
	if c.authService == nil {
		return nil, fmt.Errorf("auth service not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create gRPC request
	registerReq := &proto.RegisterRequest{
		Email:     email,
		Password:  password,
		ProjectId: projectID,
	}

	// Make gRPC call
	registerResp, err := c.authService.Register(ctx, registerReq)
	if err != nil {
		return nil, fmt.Errorf("gRPC register failed: %w", err)
	}

	// Convert gRPC response to our response format
	return &AuthResponse{
		UserID:       registerResp.GetUserId(),
		SessionToken: registerResp.GetSessionToken(),
		Email:        email,
		Valid:        true,
		ProjectIDs:   []string{projectID},
	}, nil
}

// ValidateSession checks if a session token is valid
func (c *AuthClient) ValidateSession(sessionToken string) (*AuthResponse, error) {
	if c.authService == nil {
		return nil, fmt.Errorf("auth service not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create gRPC request
	validateReq := &proto.ValidateSessionRequest{
		SessionToken: sessionToken,
	}

	// Make gRPC call
	validateResp, err := c.authService.ValidateSession(ctx, validateReq)
	if err != nil {
		return nil, fmt.Errorf("gRPC validate session failed: %w", err)
	}

	// Convert gRPC response to our response format
	return &AuthResponse{
		UserID:     validateResp.GetUserId(),
		Valid:      validateResp.GetValid(),
		ProjectIDs: validateResp.GetProjectIds(),
	}, nil
}

// HealthCheck checks if the auth service is available
func (c *AuthClient) HealthCheck() (*AuthResponse, error) {
	if c.authService == nil {
		return nil, fmt.Errorf("auth service not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create gRPC request
	healthReq := &proto.HealthCheckRequest{
		Service: "go-templ-htmx-ex",
	}

	// Make gRPC call
	healthResp, err := c.authService.HealthCheck(ctx, healthReq)
	if err != nil {
		return nil, fmt.Errorf("gRPC health check failed: %w", err)
	}

	// Convert gRPC response to our response format
	return &AuthResponse{
		UserID: "health-check",
		Valid:  healthResp.GetStatus() == "healthy",
	}, nil
}

// Close closes the gRPC connection
func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
