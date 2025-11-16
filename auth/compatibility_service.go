package auth

import (
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// CompatibilityService provides backward compatibility methods
type CompatibilityService struct {
	Service *Service
}

// NewCompatibilityService creates a new compatibility service
func NewCompatibilityService(service *Service) *CompatibilityService {
	return &CompatibilityService{Service: service}
}

// CallAuthService makes a request to the auth microservice
func (cs *CompatibilityService) CallAuthService(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	req, err := cs.Service.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	_, bodyBytes, err := cs.Service.http.Do(req)
	if err != nil {
		return nil, err
	}

	authResp, err := cs.Service.parser.ParseAuthResponse(bodyBytes)
	if err != nil {
		return nil, err
	}

	return authResp, cs.Service.parser.ValidateResponseSuccess(authResp)
}

// Logout logs out a user
func (cs *CompatibilityService) Logout(token string) error {
	fmt.Printf("User logged out with token: %s\n", token)
	return nil
}

// CreateSession exchanges OAuth authorization code for session creation
func (cs *CompatibilityService) CreateSession(code string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/auth/session/create", cs.Service.config.AuthServiceURL)
	params := map[string]string{"code": code}

	req, err := cs.Service.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	_, bodyBytes, err := cs.Service.http.Do(req)
	if err != nil {
		return nil, err
	}

	return cs.Service.parser.ParseGenericResponse(bodyBytes)
}

// ValidateToken validates a token (alias for ValidateSession)
func (cs *CompatibilityService) ValidateToken(token string) (*models.AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/refresh", cs.Service.config.AuthServiceURL)
	params := map[string]string{"session_id": token}
	return cs.CallAuthService(endpoint, params)
}