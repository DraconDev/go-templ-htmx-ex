package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// Service handles session management with the auth microservice
type Service struct {
	config    *config.Config
	http      *HTTPClient
	timeout   time.Duration
	logger    *log.Logger
	mu        sync.RWMutex
	sessions  map[string]*SessionInfo
	metrics   *ServiceMetrics
}

// SessionInfo represents session information
type SessionInfo struct {
	SessionID   string    `json:"session_id"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	LastAccess  time.Time `json:"last_access"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	IsActive    bool      `json:"is_active"`
}

// ServiceMetrics tracks service metrics
type ServiceMetrics struct {
	TotalRequests    int64     `json:"total_requests"`
	ActiveSessions   int64     `json:"active_sessions"`
	FailedRequests   int64     `json:"failed_requests"`
	LastRequest      time.Time `json:"last_request"`
	mu               sync.RWMutex
}

// NewService creates a new auth service instance
func NewService(cfg *config.Config) *Service {
	return &Service{
		config:  cfg,
		http:    NewHTTPClient(cfg.AuthServiceURL, 10*time.Second),
		timeout: 10 * time.Second,
	}
}

// CreateSession exchanges OAuth authorization code for session_id and user info
func (s *Service) CreateSession(auth_code string) (map[string]interface{}, error) {
	return s.callAuthServiceGeneric("/auth/session/create", map[string]string{
		"auth_code": auth_code,
	})
}

// RefreshSession refreshes an existing session_id
func (s *Service) RefreshSession(session_id string) (*models.AuthResponse, error) {
	return s.callAuthService("/auth/session/refresh", map[string]string{
		"session_id": session_id,
	})
}

// GetUserInfo retrieves user information using session_id
func (s *Service) GetUserInfo(session_id string) (*models.AuthResponse, error) {
	return s.callAuthService("/auth/userinfo", map[string]string{
		"session_id": session_id,
	})
}

// Logout logs out a user using session_id
func (s *Service) Logout(session_id string) error {
	fmt.Printf("User logged out with session_id: %s\n", session_id)
	return nil
}

// callAuthService makes a request to the auth microservice
func (s *Service) callAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.http.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	_, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}

	var authResp models.AuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

// callAuthServiceGeneric makes a request and returns generic response
func (s *Service) callAuthServiceGeneric(endpoint string, params map[string]string) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.http.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	_, bodyBytes, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}
