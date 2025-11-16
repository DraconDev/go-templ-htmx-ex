package session

import (
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/auth/http"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// Validator handles session validation operations
type Validator struct {
	parser  *http.Parser
	builder *http.RequestBuilder
	client  *http.Client
	baseURL string
}

// NewValidator creates a new session validator
func NewValidator(client *http.Client, builder *http.RequestBuilder, parser *http.Parser, baseURL string) *Validator {
	return &Validator{
		parser:  parser,
		builder: builder,
		client:  client,
		baseURL: baseURL,
	}
}

// ValidateSession validates a session token
func (v *Validator) ValidateSession(sessionID string) (*models.AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/session/refresh", v.baseURL)
	params := map[string]string{
		"session_id": sessionID,
	}
	
	req, err := v.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	_, bodyBytes, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}

	authResp, err := v.parser.ParseAuthResponse(bodyBytes)
	if err != nil {
		return nil, err
	}

	return authResp, v.parser.ValidateResponseSuccess(authResp)
}