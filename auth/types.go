package auth

import (
	"fmt"
	"time"
)

// AuthServiceError represents errors from the auth service
type AuthServiceError struct {
	Code    string
	Message string
	Err     error
}

func (e *AuthServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("auth service error: %s - %v", e.Message, e.Err)
	}
	return fmt.Sprintf("auth service error: %s", e.Message)
}

func (e *AuthServiceError) Unwrap() error {
	return e.Err
}

// NewAuthServiceError creates a new auth service error
func NewAuthServiceError(code, message string, err error) *AuthServiceError {
	return &AuthServiceError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// HTTPClientConfig represents HTTP client configuration
type HTTPClientConfig struct {
	Timeout       time.Duration
	AuthSecret    string
	AuthServiceURL string
}