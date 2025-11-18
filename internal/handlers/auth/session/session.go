package session

import (
	"github.com/DraconDev/go-templ-htmx-ex/internal/services"
	"github.com/DraconDev/go-templ-htmx-ex/internal/utils/config"
)

type SessionHandler struct {
	Config      *config.Config
	AuthService *services.AuthService
}

func NewSessionHandler(config *config.Config) *SessionHandler {
	return &SessionHandler{
		Config:      config,
		AuthService: services.NewAuthService(config),
	}
}
