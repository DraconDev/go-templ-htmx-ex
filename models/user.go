package models

// AuthResponse represents the response from the auth service
type AuthResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Token    string `json:"token,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Picture  string `json:"picture,omitempty"`
	Error    string `json:"error,omitempty"`
}

// TokenExchangeResponse represents the response from exchanging auth code for session
type TokenExchangeResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	IdToken     string `json:"id_token"` // Server session token
	Error       string `json:"error,omitempty"`
}

// ExchangeCodeRequest represents a request to exchange authorization code for tokens
type ExchangeCodeRequest struct {
	AuthCode string `json:"auth_code"`
}