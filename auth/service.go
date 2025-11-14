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

// ParseJWTFromIDToken extracts user information from the id_token JWT
func (s *Service) ParseJWTFromIDToken(idToken string) (*models.JWTClaims, error) {
	fmt.Printf("🔄 AUTHSVC: Parsing JWT from id_token...\n")

	// Split JWT into parts
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format: expected 3 parts, got %d", len(parts))
	}

	// Decode the payload (middle part)
	payload, err := jwtBase64URLDecode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT payload: %v", err)
	}

	// Parse user data from JWT payload
	var claims models.JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse JWT claims: %v", err)
	}

	fmt.Printf("🔄 AUTHSVC: JWT claims extracted - Subject: %s, Name: %s, Email: %s\n",
		claims.Subject, claims.Name, claims.Email)

	return &claims, nil
}

// jwtBase64URLDecode decodes base64url encoding (needed for JWT)
func jwtBase64URLDecode(data string) ([]byte, error) {
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

// CreateSession exchanges OAuth authorization code for session creation
// This is a test function to see what /session/create returns
func (s *Service) CreateSession(code string) (interface{}, error) {
	fmt.Printf("🔄 AUTHSVC: Testing /session/create endpoint...\n")

	// Call the authentication endpoint (handles both creation and refreshing)
	fmt.Printf("🔄 AUTHSVC: Calling authenticate endpoint: %s/auth/authenticate\n", s.config.AuthServiceURL)

	client := &http.Client{Timeout: 10 * time.Second}

	// Create request with the code
	jsonData := map[string]string{"code": code}
	reqData, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/authenticate", s.config.AuthServiceURL), bytes.NewBuffer(reqData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Add auth secret if configured
	if s.config.AuthSecret != "" {
		req.Header.Set("X-Auth-Secret", s.config.AuthSecret)
	}

	fmt.Printf("🔄 AUTHSVC: Sending request...\n")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("🔄 AUTHSVC: Response status: %s\n", resp.Status)

	// Read response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("🔄 AUTHSVC: Response body: %s\n", string(bodyBytes))

	// Parse response as JSON
	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}

	fmt.Printf("🔄 AUTHSVC: Parsed response: %+v\n", response)
	return response, nil
}

// RefreshToken refreshes a token using the refresh token
func (s *Service) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	fmt.Printf("🔄 AUTHSVC: Refreshing token with refresh_token length: %d\n", len(refreshToken))

	// Special handling for refresh endpoint - it returns tokens, not user info
	return s.CallRefreshAuthService(fmt.Sprintf("%s/auth/refresh", s.config.AuthServiceURL), map[string]string{
		"refresh_token": refreshToken,
	})
}

// CallRefreshAuthService handles the special response format from refresh endpoint
func (s *Service) CallRefreshAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	fmt.Printf("🔄 AUTHSVC: Calling refresh endpoint: %s\n", endpoint)

	client := &http.Client{Timeout: 10 * time.Second}

	// Create JSON data
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	fmt.Printf("🔄 AUTHSVC: Refresh request data: %s\n", string(jsonData))

	// Add auth secret to params if not present
	if _, exists := params["secret"]; !exists && s.config.AuthSecret != "" {
		params["secret"] = s.config.AuthSecret
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Add X-Auth-Secret header for protected endpoints
	if s.config.AuthSecret != "" {
		req.Header.Set("X-Auth-Secret", s.config.AuthSecret)
	}

	fmt.Printf("🔄 AUTHSVC: Sending refresh request...\n")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("🔄 AUTHSVC: Refresh response status: %s\n", resp.Status)

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("🔄 AUTHSVC: Refresh response body: %s\n", string(bodyBytes))

	// Parse the refresh response (has id_token and refresh_token fields)
	var refreshResp struct {
		IdToken      string `json:"id_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.Unmarshal(bodyBytes, &refreshResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %v", err)
	}

	if refreshResp.IdToken == "" {
		return nil, fmt.Errorf("no id_token in refresh response")
	}

	fmt.Printf("🔄 AUTHSVC: Successfully got new tokens - id_token: %d chars, refresh_token: %d chars\n",
		len(refreshResp.IdToken), len(refreshResp.RefreshToken))

	// Parse the new id_token to extract user info
	claims, err := s.ParseJWTFromIDToken(refreshResp.IdToken)
	if err != nil {
		fmt.Printf("🔄 AUTHSVC: Failed to parse new id_token: %v\n", err)
		// Return basic response without user data
		return &models.AuthResponse{
			Success: true,
			Token:   refreshResp.IdToken,
		}, nil
	}

	return &models.AuthResponse{
		Success: true,
		Token:   refreshResp.IdToken,
		Name:    claims.Name,
		Email:   claims.Email,
		Picture: claims.Picture,
		UserID:  claims.Subject,
	}, nil
}

// ExchangeCodeForTokens exchanges OAuth authorization code for session and refresh tokens
func (s *Service) ExchangeCodeForTokens(code string) (*models.TokenExchangeResponse, error) {
	fmt.Printf("🔄 AUTHSVC: Exchanging code for tokens...\n")

	// Call the auth service directly and parse the response manually
	fmt.Printf("🔄 AUTHSVC: Calling auth service endpoint: %s/auth/token\n", s.config.AuthServiceURL)

	client := &http.Client{Timeout: 10 * time.Second}

	// Create request
	jsonData := map[string]string{"code": code}
	reqData, err := json.Marshal(jsonData)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to marshal request data",
		}, err
	}
	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/session/validate", s.config.AuthServiceURL), bytes.NewBuffer(reqData))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/token", s.config.AuthServiceURL), bytes.NewBuffer(reqData))
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to create request",
		}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Add auth secret if configured
	if s.config.AuthSecret != "" {
		req.Header.Set("X-Auth-Secret", s.config.AuthSecret)
	}

	resp, err := client.Do(req)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to call auth service: " + err.Error(),
		}, err
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to read response body",
		}, err
	}

	fmt.Printf("🔄 AUTHSVC: Response status: %s\n", resp.Status)
	fmt.Printf("🔄 AUTHSVC: Response body: %s\n", string(bodyBytes))

	// Parse the response directly as a map to extract tokens
	var respData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &respData); err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to parse response: " + err.Error(),
		}, fmt.Errorf("failed to parse auth service response: %v", err)
	}

	// Check if we have an error in the response
	if errMsg, hasError := respData["error"]; hasError {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   fmt.Sprintf("%v", errMsg),
		}, fmt.Errorf("auth service error: %v", errMsg)
	}

	// Extract id token and refresh token
	var idToken, refreshToken string
	var hasIdToken, hasRefresh bool

	if idTokenInterface, exists := respData["id_token"]; exists {
		if idTokenStr, ok := idTokenInterface.(string); ok {
			idToken = idTokenStr
			hasIdToken = true
		}
	}

	if refreshInterface, exists := respData["refresh_token"]; exists {
		if refreshStr, ok := refreshInterface.(string); ok {
			refreshToken = refreshStr
			hasRefresh = true
		}
	}

	fmt.Printf("🔄 AUTHSVC: Token extraction - IdToken: %t (%d chars), Refresh: %t (%d chars)\n",
		hasIdToken, len(idToken), hasRefresh, len(refreshToken))

	if !hasIdToken || !hasRefresh || idToken == "" || refreshToken == "" {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   fmt.Sprintf("Missing tokens - IdToken: %t, Refresh: %t", hasIdToken, hasRefresh),
		}, fmt.Errorf("missing tokens in auth service response")
	}

	fmt.Printf("🔄 AUTHSVC: Successfully extracted tokens - IdToken: %d chars, Refresh: %d chars\n",
		len(idToken), len(refreshToken))

	return &models.TokenExchangeResponse{
		Success:      true,
		IdToken:      idToken,
		RefreshToken: refreshToken,
	}, nil
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
