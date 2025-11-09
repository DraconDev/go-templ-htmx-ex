package models

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	LoggedIn bool   `json:"logged_in"`
}

// Config holds application configuration
type Config struct {
	ServerPort     string
	AuthServiceURL string
	RedirectURL    string
}

// UserSession represents a logged-in user session
type UserSession struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	LoggedIn bool   `json:"logged_in"`
}