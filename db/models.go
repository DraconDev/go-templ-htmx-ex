package db

import "time"

// User represents a user in the application
// Complements the external auth service by storing app-level user data
type User struct {
	ID        string    `json:"id" db:"id"`
	AuthID    string    `json:"auth_id" db:"auth_id"` // References auth service user ID
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Picture   string    `json:"picture" db:"picture"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"` // Admin flag for role-based access
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserPreferences represents user application preferences
// Essential for any startup to provide good user experience
type UserPreferences struct {
	ID                string    `json:"id" db:"id"`
	UserID            string    `json:"user_id" db:"user_id"`
	Theme             string    `json:"theme" db:"theme"`       // "dark", "light", "auto"
	Language          string    `json:"language" db:"language"` // "en", "es", "fr", etc.
	EmailNotifications bool      `json:"email_notifications" db:"email_notifications"`
	PushNotifications  bool      `json:"push_notifications" db:"push_notifications"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}