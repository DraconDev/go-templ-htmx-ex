package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/config"
	dbSqlc "github.com/DraconDev/go-templ-htmx-ex/db/sqlc"
	"github.com/DraconDev/go-templ-htmx-ex/templates"
)

// AdminHandler handles admin-specific operations
type AdminHandler struct {
	Config   *config.Config
	Queries  *dbSqlc.Queries // SQLC generated queries
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(config *config.Config, queries *dbSqlc.Queries) *AdminHandler {
	return &AdminHandler{
		Config:  config,
		Queries: queries,
	}
}

// AdminDashboardHandler serves the admin dashboard
func (h *AdminHandler) AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("📋 ADMIN: Admin dashboard requested\n")
	
	// Get user info using existing JWT validation logic
	userInfo := GetUserFromJWT(r)
	
	if !userInfo.LoggedIn {
		fmt.Printf("📋 ADMIN: User not logged in\n")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fmt.Printf("📋 ADMIN: User logged in: %s (%s)\n", userInfo.Name, userInfo.Email)

	// Check if this user is admin
	if !h.Config.IsAdmin(userInfo.Email) {
		fmt.Printf("📋 ACCESS DENIED: User %s is not admin (admin email: %s)\n",
			userInfo.Email, h.Config.AdminEmail)
		http.Error(w, "Access denied: Admin privileges required", http.StatusForbidden)
		return
	}

	fmt.Printf("📋 ADMIN: Access granted for admin %s\n", userInfo.Email)

	w.Header().Set("Content-Type", "text/html")
	component := templates.Layout("Admin Dashboard", templates.NavigationLoggedIn(userInfo), templates.AdminDashboardContent(userInfo))
	component.Render(r.Context(), w)
}

// GetUsersHandler returns a list of users from the database
func (h *AdminHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if h.Database == nil {
		// Fallback to enhanced mock data if no database connection
		fmt.Println("🔍 Using mock user data - no database connection")
		users := []map[string]interface{}{
			{
				"id":        1,
				"email":     "john.doe@example.com",
				"name":      "John Doe",
				"picture":   "https://ui-avatars.com/api/?name=John+Doe&background=3B82F6&color=fff&size=40",
				"role":      "user",
				"status":    "active",
				"lastLogin": "2025-11-11T20:45:00Z",
				"createdAt": "2025-11-10T14:30:00Z",
			},
			{
				"id":        2,
				"email":     "dracsharp@gmail.com",
				"name":      "Platform Admin",
				"picture":   "https://ui-avatars.com/api/?name=Admin&background=EF4444&color=fff&size=40",
				"role":      "admin",
				"status":    "active",
				"lastLogin": "2025-11-11T21:00:00Z",
				"createdAt": "2025-11-08T12:00:00Z",
			},
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": users,
			"total": len(users),
			"active": 2,
			"inactive": 0,
		})
		return
	}

	// Try to get real users from database
	userRepo := db.NewUserRepository(h.Database)
	
	// Try to get the admin user from database first
	adminUser, err := userRepo.GetUserByEmail(h.Config.AdminEmail)
	if err == nil && adminUser != nil {
		fmt.Printf("✅ Found admin user in database: %s\n", adminUser.Email)
		
		// Convert admin user to response format
		userMaps := []map[string]interface{}{
			{
				"id":        adminUser.ID,
				"email":     adminUser.Email,
				"name":      adminUser.Name,
				"picture":   adminUser.Picture,
				"role":      "admin",
				"status":    "active",
				"lastLogin": adminUser.CreatedAt.Format("2006-01-02T15:04:05Z"),
				"createdAt": adminUser.CreatedAt.Format("2006-01-02T15:04:05Z"),
			},
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": userMaps,
			"total": len(userMaps),
			"active": len(userMaps),
			"inactive": 0,
		})
	} else {
		fmt.Printf("🔍 Admin user not found in database or DB query failed: %v\n", err)
		// Fallback to enhanced mock data including the admin email
		users := []map[string]interface{}{
			{
				"id":        1,
				"email":     "john.doe@example.com",
				"name":      "John Doe",
				"picture":   "https://ui-avatars.com/api/?name=John+Doe&background=3B82F6&color=fff&size=40",
				"role":      "user",
				"status":    "active",
				"lastLogin": "2025-11-11T20:45:00Z",
				"createdAt": "2025-11-10T14:30:00Z",
			},
			{
				"id":        2,
				"email":     h.Config.AdminEmail,
				"name":      "Platform Admin",
				"picture":   "https://ui-avatars.com/api/?name=Admin&background=EF4444&color=fff&size=40",
				"role":      "admin",
				"status":    "active",
				"lastLogin": "2025-11-11T21:01:00Z",
				"createdAt": "2025-11-08T12:00:00Z",
			},
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": users,
			"total": len(users),
			"active": len(users),
			"inactive": 0,
		})
	}
}

// GetAnalyticsHandler returns analytics data (stub for now)
func (h *AdminHandler) GetAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Mock analytics data
	analytics := map[string]interface{}{
		"total_users":      127,
		"active_sessions":  23,
		"signups_today":   5,
		"signups_this_week": 34,
		"system_health":    "operational",
	}

	json.NewEncoder(w).Encode(analytics)
}

// GetSettingsHandler returns system settings (stub for now)
func (h *AdminHandler) GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Mock settings data
	settings := map[string]interface{}{
		"maintenance_mode": false,
		"registration_enabled": true,
		"max_users": 1000,
		"session_timeout": 3600,
	}

	json.NewEncoder(w).Encode(settings)
}

// GetLogsHandler returns system logs (stub for now)
func (h *AdminHandler) GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Mock log data
	logs := []map[string]interface{}{
		{
			"timestamp": "2025-11-11T02:58:00Z",
			"level":     "INFO",
			"message":   "User login successful",
			"user":      "john@example.com",
		},
		{
			"timestamp": "2025-11-11T02:57:00Z",
			"level":     "WARN",
			"message":   "High memory usage detected",
			"user":      "system",
		},
		{
			"timestamp": "2025-11-11T02:56:00Z",
			"level":     "ERROR",
			"message":   "Database connection failed",
			"user":      "system",
		},
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs": logs,
		"total": len(logs),
	})
}