
package services

import (
	"encoding/json"
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/internal/models"
	"github.com/DraconDev/go-templ-htmx-ex/internal/utils/config"
)

type AuthService struct {
	config    *config.Config
	authClient *AuthClient
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		config:    cfg,
		authClient: NewAuthClient(cfg.AuthServiceURL),
	}
}

func (s *AuthService) CreateSession(auth_code string) (map[string]interface{},
