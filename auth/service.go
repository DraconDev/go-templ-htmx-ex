package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// Service handles session management with the auth microservice
type Service struct {
	config  *config.Config
	http    *HTTPClient
	timeout time.Duration
	logger  *log.Logger
}

// AuthError represents authentication errors
type AuthError struct {
	Code    string
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

// NewService creates a new auth service instance
func NewService(cfg *config.Config) *Service {
	return &Service{
		config:  cfg,
		http:    NewHTTPClient(cfg.AuthServiceURL, 10*time.Second),
		timeout: 10 * time.Second,
		logger:  log.New(log.Writer(), "[auth-service] ", log.LstdFlags),
	}
}

// makeRequest handles the common HTTP request logic
func (s *Service) makeRequest(endpoint string, params map[string]string) (*http.Response, []byte, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, nil, &AuthError{
			Code:    "SERIALIZATION_ERROR",
			Message: fmt.Sprintf("Failed to serialize request: %v", err),
		}
	}

	req, err := http.NewRequest("POST", s.http.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, &AuthError{
			Code:    "REQUEST_ERROR",
			Message: fmt.Sprintf("Failed to create request: %v", err),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	s.logger.Printf("Making request to: %s", endpoint)

	return s.http.Do(req)
}

// callAuthService makes a request to the auth microservice
func (s *Service) callAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	resp, bodyBytes, err := s.makeRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &AuthError{
			Code:    "HTTP_ERROR",
			Message: fmt.Sprintf("Auth service returned status %d", resp.StatusCode),
		}
	}

	var authResp models.AuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		return nil, &AuthError{
			Code:    "DESERIALIZATION_ERROR",
			Message: fmt.Sprintf("Failed to unmarshal response: %v", err),
		}
	}

	if !authResp.Success {
		return nil, &AuthError{
			Code:    "AUTH_SERVICE_ERROR",
			Message: authResp.Error,
		}
	}

	return &authResp, nil
}

// callAuthServiceGeneric makes a request and returns generic response
func (s *Service) callAuthServiceGeneric(endpoint string, params map[string]string) (map[string]interface{}, error) {
	resp, bodyBytes, err := s.makeRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &AuthError{
			Code:    "HTTP_ERROR",
			Message: fmt.Sprintf("Auth service returned status %d", resp.StatusCode),
		}
	}

	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, &AuthError{
			Code:    "DESERIALIZATION_ERROR",
			Message: fmt.Sprintf("Failed to unmarshal response: %v", err),
		}
	}

	return response, nil
}

// HealthCheck verifies the auth service is responsive
func (s *Service) HealthCheck() error {
	s.logger.Println("Performing health check")
	
	_, err := s.callAuthServiceGeneric("/health", map[string]string{})
	if err != nil {
		s.logger.Printf("Health check failed: %v", err)
		return err
	}
	
	s.logger.Println("Health check passed")
	return nil
}
