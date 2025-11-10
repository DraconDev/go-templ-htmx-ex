package models

// UserSession represents a logged-in user session
type UserSession struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	LoggedIn bool   `json:"logged_in"`
}

// AuthResponse represents the response from the auth service
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
	Picture string `json:"picture,omitempty"`
	Error   string `json:"error,omitempty"`
}

// SetSessionRequest represents a request to set a session token
type SetSessionRequest struct {
	Token string `json:"token"`
}