package admin

import (
	"github.com/DraconDev/go-templ-htmx-ex/config"
	dbSqlc "github.com/DraconDev/go-templ-htmx-ex/db/sqlc"
)

// AdminHandler handles admin-specific operations
type AdminHandler struct {
	Config  *config.Config
	Queries *dbSqlc.Queries // SQLC generated queries
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(config *config.Config, queries *dbSqlc.Queries) *AdminHandler {
	return &AdminHandler{
		Config:  config,
		Queries: queries,
	}
}