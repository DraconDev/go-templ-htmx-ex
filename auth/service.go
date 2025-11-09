package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// Service handles communication with the auth microservice
type Service struct {
	config *config.Config
}

// NewService creates a new auth service instance
func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

// getStringFromMap safely gets a string from a map
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// CallAuthService makes a request to the auth microservice
func (s *Service) CallAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	fmt.Printf("🔐 AUTHSVC: === CallAuthService STARTED ===\n")
	fmt.Printf("🔐 AUTHSVC: Endpoint: %s\n", endpoint)
	fmt.Printf("🔐 AUTHSVC: Params: %v\n", params)
	
	client := &http.Client{Timeout: 10 * time.Second}

	// Create JSON data
	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Printf("🔐 AUTHSVC: JSON marshaling failed: %v\n", err)
		return nil, err
	}
	
	fmt.Printf("🔐 AUTHSVC: JSON data: %s\n", string(jsonData))

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("🔐 AUTHSVC: Request creation failed: %v\n", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("🔐 AUTHSVC: Sending request to auth service...\n")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("🔐 AUTHSVC: Request failed: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	
	fmt.Printf("🔐 AUTHSVC: Response status: %s\n", resp.Status)

	// Read the response body first
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("🔐 AUTHSVC: Failed to read response body: %v\n", err)
		return nil, err
	}
	
	fmt.Printf("🔐 AUTHSVC: Response body: %s\n", string(bodyBytes))

	// Try to decode as AuthResponse first
	var authResp models.AuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err == nil && authResp.Success {
		fmt.Printf("🔐 AUTHSVC: Parsed as AuthResponse - Success: %v\n", authResp.Success)
		return &authResp, nil
	}

	fmt.Printf("🔐 AUTHSVC: Failed to parse as AuthResponse, trying JWT payload...\n")
	// If that fails, try to decode as JWT payload and convert
	var jwtPayload map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jwtPayload); err != nil {
		fmt.Printf("🔐 AUTHSVC: Failed to parse JWT payload: %v\n", err)
		return nil, err
	}

	// Convert JWT payload to AuthResponse format
	result := &models.AuthResponse{
		Success: true,
		Name:    getStringFromMap(jwtPayload, "name"),
		Email:   getStringFromMap(jwtPayload, "email"),
		Picture: getStringFromMap(jwtPayload, "picture"),
		UserID:  getStringFromMap(jwtPayload, "sub"),
	}
	
	fmt.Printf("🔐 AUTHSVC: Converted JWT to AuthResponse - Name: %s, Email: %s\n", result.Name, result.Email)
	fmt.Printf("🔐 AUTHSVC: === CallAuthService COMPLETED ===\n")
	
	return result, nil
}

// ValidateSession validates a session token
func (s *Service) ValidateSession(token string) (*models.AuthResponse, error) {
	return s.CallAuthService(fmt.Sprintf("%s/auth/validate", s.config.AuthServiceURL), map[string]string{
		"token": token,
	})
}

// GetUserInfo retrieves user information from auth service
func (s *Service) GetUserInfo(token string) (*models.AuthResponse, error) {
	return s.CallAuthService(fmt.Sprintf("%s/auth/userinfo", s.config.AuthServiceURL), map[string]string{
		"token": token,
	})
}

// Logout logs out a user (通知auth service)
func (s *Service) Logout(token string) error {
	// Since this is a JWT-based system, we can just log it
	// In a more complex system, you might want to blacklist the token
	log.Printf("User logged out with token: %s", token)
	return nil
}

// ValidateUser validates a user token (alias for GetUserInfo)
func (s *Service) ValidateUser(token string) (*models.AuthResponse, error) {
	return s.GetUserInfo(token)
}

// ValidateToken validates a token (alias for ValidateSession)
func (s *Service) ValidateToken(token string) (*models.AuthResponse, error) {
	return s.ValidateSession(token)
}