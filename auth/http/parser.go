package http

import (
	"encoding/json"
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// ResponseParser handles parsing responses from the auth service
type ResponseParser struct{}

// NewResponseParser creates a new response parser
func NewResponseParser() *ResponseParser {
	return &ResponseParser{}
}

// ParseAuthResponse parses a response as AuthResponse
func (rp *ResponseParser) ParseAuthResponse(bodyBytes []byte) (*models.AuthResponse, error) {
	var authResp models.AuthResponse
	if err := json.Unmarshal(bodyBytes, &authResp); err != nil {
		return nil, NewAuthServiceError("PARSE_FAILED", "Failed to parse auth response", err)
	}

	return &authResp, nil
}

// ParseTokenExchangeResponse parses a response as TokenExchangeResponse
func (rp *ResponseParser) ParseTokenExchangeResponse(bodyBytes []byte) (*models.TokenExchangeResponse, error) {
	var respData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &respData); err != nil {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Failed to parse response: " + err.Error(),
		}, NewAuthServiceError("PARSE_FAILED", "Failed to parse token exchange response", err)
	}

	// Check for errors in response
	if errMsg, hasError := respData["error"]; hasError {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   fmt.Sprintf("%v", errMsg),
		}, NewAuthServiceError("AUTH_SERVICE_ERROR", fmt.Sprintf("Auth service error: %v", errMsg), nil)
	}

	// Extract session_id
	var sessionID string
	if sessionInterface, exists := respData["session_id"]; exists {
		if sessionStr, ok := sessionInterface.(string); ok {
			sessionID = sessionStr
		}
	}

	if sessionID == "" {
		return &models.TokenExchangeResponse{
			Success: false,
			Error:   "Missing session_id in auth service response",
		}, NewAuthServiceError("MISSING_SESSION_ID", "Missing session_id in response", nil)
	}

	return &models.TokenExchangeResponse{
		Success: true,
		IdToken: sessionID,
	}, nil
}

// ParseGenericResponse parses a response as generic map
func (rp *ResponseParser) ParseGenericResponse(bodyBytes []byte) (map[string]interface{}, error) {
	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, NewAuthServiceError("PARSE_FAILED", "Failed to parse generic response", err)
	}

	return response, nil
}

// ValidateResponseSuccess checks if the response indicates success
func (rp *ResponseParser) ValidateResponseSuccess(authResp *models.AuthResponse) error {
	if !authResp.Success {
		return NewAuthServiceError("AUTH_FAILED", "Authentication failed", nil)
	}
	return nil
}