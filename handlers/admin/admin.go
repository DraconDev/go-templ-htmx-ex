package admin

import (
	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/services"
	dbSqlc "github.com/DraconDev/go-templ-htmx-ex/db/sqlc"
)

// AdminHandler handles admin-specific operations
type AdminHandler struct {
	Config      *config.Config
	UserService *services.UserService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(config *config.Config, queries *dbSqlc.Queries) *AdminHandler {
	return &AdminHandler{
		Config:      config,
		UserService: services.NewUserService(queries),
	}
}