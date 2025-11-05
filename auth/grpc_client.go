package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"star/auth/cerberus"
)

// GRPCAuthClient handles communication with the Cerberus authentication service via gRPC
type GRPCAuthClient struct {
	client    auth_cerberus.AuthServiceClient
	conn      *grpc.ClientConn
	ctx       context.Context
	cancel    context.CancelFunc
	baseURL   string
}

// GRPCAuthResponse represents the response from gRPC auth service
type GRPCAuthResponse struct {
	UserID       string   `json:"user_id"`
	SessionToken string   `json:"session_token"`
	Email        string   `json:"email,omitempty"`
	Valid        bool     `json:"valid,omitempty"`
	ProjectIDs   []string `json:"project_ids,omitempty"`
	Error        string   `json:"error,omitempty"`
	Status       string   `json:"status,omitempty"`
	Message      string   `json:"message,omitempty"`
}

// NewGRPCAuthClient creates a new gRPC authentication client
func NewGRPCAuthClient(baseURL string) (*GRPCAuthClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	
	// Set up gRPC connection with insecure credentials for now
	// In production, you would use proper TLS certificates
	conn, err := grpc.NewClient(
		baseURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	client := auth_cerberus.NewAuthServiceClient(conn)
	
	return &GRPCAuthClient{
		client:  client,
		conn:    conn,
		ctx:     ctx,
		cancel:  cancel,
		baseURL: baseURL,
	}, nil
}

// Close closes the gRPC connection
func (c *GRPCAuthClient) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Login attempts to authenticate a user and returns JWT token
func (c *GRPCAuthClient) Login(email, password string) (*GRPCAuthResponse, error) {
	req := &auth_cerberus.LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.client.Login(c.ctx, req)
	if err != nil {
		// Handle gRPC errors
		st, ok := status.FromError(err)
		if ok {
			return &GRPCAuthResponse{
				Error: fmt.Sprintf("gRPC error (%s): %s", st.Code(), st.Message()),
			}, err
		}
		return &GRPCAuthResponse{
			Error: fmt.Sprintf("login failed: %v", err),
		}, err
	}

	return &GRPCAuthResponse{
		UserID:       resp.GetUserId(),
		SessionToken: resp.GetSessionToken(),
		Status:       "success",
	}, nil
}

// Register creates a new user account
func (c *GRPCAuthClient) Register(email, password, projectID string) (*GRPCAuthResponse, error) {
	req := &auth_cerberus.RegisterRequest{
		Email:     email,
		Password:  password,
		ProjectId: projectID,
	}

	resp, err := c.client.Register(c.ctx, req)
	if err != nil {
		// Handle gRPC errors
		st, ok := status.FromError(err)
		if ok {
			return &GRPCAuthResponse{
				Error: fmt.Sprintf("gRPC error (%s): %s", st.Code(), st.Message()),
			}, err
		}
		return &GRPCAuthResponse{
			Error: fmt.Sprintf("registration failed: %v", err),
		}, err
	}

	return &GRPCAuthResponse{
		UserID:       resp.GetUserId(),
		SessionToken: resp.GetSessionToken(),
		Status:       "success",
	}, nil
}

// ValidateSession checks if a session token is valid
func (c *GRPCAuthClient) ValidateSession(sessionToken string) (*GRPCAuthResponse, error) {
	req := &auth_cerberus.ValidateSessionRequest{
		SessionToken: sessionToken,
	}

	resp, err := c.client.ValidateSession(c.ctx, req)
	if err != nil {
		// Handle gRPC errors
		st, ok := status.FromError(err)
		if ok {
			return &GRPCAuthResponse{
				Error: fmt.Sprintf("gRPC error (%s): %s", st.Code(), st.Message()),
			}, err
		}
		return &GRPCAuthResponse{
			Error: fmt.Sprintf("session validation failed: %v", err),
		}, err
	}

	return &GRPCAuthResponse{
		UserID:     resp.GetUserId(),
		Valid:      resp.GetValid(),
		ProjectIDs: resp.GetProjectIds(),
		Status:     "validated",
	}, nil
}

// GetUserDetails retrieves user details by user ID
func (c *GRPCAuthClient) GetUserDetails(userID string) (*GRPCAuthResponse, error) {
	req := &auth_cerberus.GetUserDetailsRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUserDetails(c.ctx, req)
	if err != nil {
		// Handle gRPC errors
		st, ok := status.FromError(err)
		if ok {
			return &GRPCAuthResponse{
				Error: fmt.Sprintf("gRPC error (%s): %s", st.Code(), st.Message()),
			}, err
		}
		return &GRPCAuthResponse{
			Error: fmt.Sprintf("get user details failed: %v", err),
		}, err
	}

	return &GRPCAuthResponse{
		UserID: resp.GetUserId(),
		Email:  resp.GetEmail(),
		Status: "success",
	}, nil
}

// HealthCheck checks if the auth service is available
func (c *GRPCAuthClient) HealthCheck() (*GRPCAuthResponse, error) {
	req := &auth_cerberus.HealthCheckRequest{
		Service: "auth-client",
	}

	resp, err := c.client.HealthCheck(c.ctx, req)
	if err != nil {
		// Handle gRPC errors
		st, ok := status.FromError(err)
		if ok {
			return &GRPCAuthResponse{
				Error: fmt.Sprintf("gRPC error (%s): %s", st.Code(), st.Message()),
			}, err
		}
		return &GRPCAuthResponse{
			Error: fmt.Sprintf("health check failed: %v", err),
		}, err
	}

	return &GRPCAuthResponse{
		Status:  resp.GetStatus(),
		Message: resp.GetMessage(),
	}, nil
}

// IsHealthy checks if the gRPC connection is healthy
func (c *GRPCAuthClient) IsHealthy() bool {
	// Simple health check - try to establish connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Attempt a simple health check call
	_, err := c.HealthCheck()
	if err != nil {
		log.Printf("gRPC client health check failed: %v", err)
		return false
	}
	
	// Check if connection state is ready
	state := c.conn.GetState()
	return state == grpc.Ready
}
