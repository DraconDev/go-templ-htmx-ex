package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/config"
	"github.com/DraconDev/go-templ-htmx-ex/db"
	"github.com/DraconDev/go-templ-htmx-ex/templates"
)

// AdminHandler handles admin-specific operations
type AdminHandler struct {
	Config *config.Config
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(config *config.Config) *AdminHandler {
	return &AdminHandler{
		Config: config,
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
	
	// Get database connection (simplified - in production you'd use proper DI)
	dbConn := db.GetConnection()
	if dbConn == nil {
		// Fallback to mock data if no database
		users := []map[string]interface{}{
			{
				"id":     1,
				"email":  "john@example.com",
				"name":   "John Doe",
				"picture": "https://via.placeholder.com/40",
				"role":    "user",
				"status":  "active",
			},
			{
				"id":     2,
				"email":  "admin@example.com",
				"name":   "Admin User",
				"picture": "https://via.placeholder.com/40",
				"role":    "admin",
				"status":  "active",
			},
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": users,
			"total": len(users),
		})
		return
	}

	// Try to get real users from database
	userRepo := db.NewUserRepository(dbConn)
	users, err := userRepo.GetAllUsers()
	if err != nil {
		// Fallback to mock data if database query fails
		users := []map[string]interface{}{
			{
				"id":     1,
				"email":  "john@example.com",
				"name":   "John Doe",
				"picture": "https://via.placeholder.com/40",
				"role":    "user",
				"status":  "active",
			},
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": users,
			"total": len(users),
		})
		return
	}

	// Convert database users to response format
	userMaps := make([]map[string]interface{}, len(users))
	for i, user := range users {
		userMaps[i] = map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"name":    user.Name,
			"picture": user.Picture,
			"role":    "user", // Default role
			"status":  "active",
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": userMaps,
		"total": len(userMaps),
	})
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