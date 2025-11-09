package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Service handles communication with the authentication microservice
type Service struct {
	BaseURL string
	Client  *http.Client
}

// NewService creates a new auth service instance
func NewService(baseURL string) *Service {
	return &Service{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// UserResponse represents the response from the auth service
type UserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
	Picture string `json:"picture,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ValidateUser validates a user token and returns user information
func (s *Service) ValidateUser(token string) (*UserResponse, error) {
	jsonData, err := json.Marshal(map[string]string{"token": token})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/auth/userinfo", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body first
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Try to decode as UserResponse first
	var userResp UserResponse
	if err := json.Unmarshal(bodyBytes, &userResp); err == nil && userResp.Success {
		return &userResp, nil
	}

	// If that fails, try to decode as JWT payload and convert
	var jwtPayload map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jwtPayload); err != nil {
		return nil, err
	}

	// Convert JWT payload to UserResponse format
	return &UserResponse{
		Success: true,
		Name:    getStringFromMap(jwtPayload, "name"),
		Email:   getStringFromMap(jwtPayload, "email"),
		Picture: getStringFromMap(jwtPayload, "picture"),
		UserID:  getStringFromMap(jwtPayload, "sub"),
	}, nil
}

// ValidateToken validates a JWT token with the auth service
func (s *Service) ValidateToken(token string) (*UserResponse, error) {
	return s.ValidateUser(token)
}

// Logout logs out a user (notifies the auth service)
func (s *Service) Logout(token string) error {
	jsonData, err := json.Marshal(map[string]string{"token": token})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/auth/logout", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logout failed with status: %d", resp.StatusCode)
	}

	return nil
}

// getStringFromMap safely gets string from map
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}