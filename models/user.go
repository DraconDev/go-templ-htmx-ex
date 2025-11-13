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
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserID       string `json:"user_id,omitempty"`
	Email        string `json:"email,omitempty"`
	Name         string `json:"name,omitempty"`
	Picture      string `json:"picture,omitempty"`
	Error        string `json:"error,omitempty"`
}

// TokenExchangeResponse represents the response from exchanging auth code for tokens
type TokenExchangeResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	IdToken      string `json:"id_token"`      // The JWT (instead of session_token)
	RefreshToken string `json:"refresh_token"` // The refresh code
	Error        string `json:"error,omitempty"`
}

// JWTClaims represents the standard OpenID Connect claims in a JWT
type JWTClaims struct {
	Subject string `json:"sub"`    // User ID
	Name    string `json:"name"`   // Full name
	Email   string `json:"email"`  // Email address
	Picture string `json:"picture"` // Avatar URL
	Issuer  string `json:"iss"`    // Issuer (auth service)
	Audience string `json:"aud"`   // Audience
	Expires int64  `json:"exp"`    // Expiration time
	Issued  int64  `json:"iat"`    // Issued at
}

// ExchangeCodeRequest represents a request to exchange authorization code for tokens
type ExchangeCodeRequest struct {
	Code string `json:"code"`
}