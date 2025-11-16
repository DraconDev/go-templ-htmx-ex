package auth

import (
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// UserInfoService handles user information operations
type UserInfoService struct {
	Service *Service
}

// NewUserInfoService creates a new user info service
func NewUserInfoService(service *Service) *UserInfoService {
	return &UserInfoService{Service: service}
}

// GetUserInfo retrieves user information from auth service
func (uis *UserInfoService) GetUserInfo(token string) (*models.AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/userinfo", uis.Service.config.AuthServiceURL)
	params := map[string]string{
		"token": token,
	}
	return uis.Service.CallAuthService(endpoint, params)
}

// ValidateUser validates a user token (alias for GetUserInfo)
func (uis *UserInfoService) ValidateUser(token string) (*models.AuthResponse, error) {
	return uis.GetUserInfo(token)
}