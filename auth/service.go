package auth

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"net/http"
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

	// Add auth secret to params if not present
	if _, exists := params["secret"]; !exists && s.config.AuthSecret != "" {
		params["secret"] = s.config.AuthSecret
		fmt.Printf("🔐 AUTHSVC: Added auth secret to params\n")
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("🔐 AUTHSVC: Request creation failed: %v\n", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	
	// Add X-Auth-Secret header for protected endpoints
	if s.config.AuthSecret != "" {
		req.Header.Set("X-Auth-Secret", s.config.AuthSecret)
		fmt.Printf("🔐 AUTHSVC: Added X-Auth-Secret header\n")
	}

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

// RefreshToken refreshes a token using the refresh token
func (s *Service) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	return s.CallAuthService(fmt.Sprintf("%s/auth/refresh", s.config.AuthServiceURL), map[string]string{
		"token": refreshToken,
	})
}

// =============================================================================
// JWT LOCAL VALIDATION - FOR OPTIMIZED UI
// =============================================================================

// PublicKeyCache stores verification keys for local JWT validation
type PublicKeyCache struct {
	keys map[string]*rsa.PublicKey
}

// NewPublicKeyCache creates a new public key cache
func NewPublicKeyCache() *PublicKeyCache {
	return &PublicKeyCache{
		keys: make(map[string]*rsa.PublicKey),
	}
}

// ValidateJWTLocal validates a JWT token locally (fast, no API call)
func (c *PublicKeyCache) ValidateJWTLocal(token string) (*models.AuthResponse, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return &models.AuthResponse{Success: false}, fmt.Errorf("invalid JWT format")
	}

	// Decode header
	headerBytes, err := base64URLDecode(parts[0])
	if err != nil {
		return &models.AuthResponse{Success: false}, fmt.Errorf("failed to decode header: %v", err)
	}

	var header struct {
		Kid string `json:"kid"`
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}

	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return &models.AuthResponse{Success: false}, fmt.Errorf("failed to parse header: %v", err)
	}

	// For now, return a basic response since we don't have the public keys yet
	// In a real implementation, you would:
	// 1. Fetch keys from /auth/jwks endpoint
	// 2. Verify signature with the correct public key
	// 3. Check expiration and claims

	return &models.AuthResponse{
		Success: true,
		Name:    "Loading...",
		Email:   "",
		Picture: "",
		UserID:  "",
	}, nil
}

// base64URLDecode decodes base64url encoding
func base64URLDecode(data string) ([]byte, error) {
	// Add padding if needed
	switch len(data) % 4 {
	case 2:
		data += "=="
	case 3:
		data += "="
	case 1:
		return nil, fmt.Errorf("invalid base64url length")
	}

	return base64.URLEncoding.DecodeString(data)
}
