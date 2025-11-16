package auth

import (
	"fmt"

	"github.com/DraconDev/go-templ-htmx-ex/models"
)

// UserServiceImpl implements user-related operations
type UserServiceImpl struct {
	http    *HTTPClient
	builder *RequestBuilder
	parser  *ResponseParser
	baseURL string
}

// NewUserServiceImpl creates a new user service
func NewUserServiceImpl(httpClient *HTTPClient, builder *RequestBuilder, parser *ResponseParser, baseURL string) *UserServiceImpl {
	return &UserServiceImpl{
		http:    httpClient,
		builder: builder,
		parser:  parser,
		baseURL: baseURL,
	}
}

// GetUserInfo retrieves user information from auth service
func (u *UserServiceImpl) GetUserInfo(token string) (*models.AuthResponse, error) {
	endpoint := fmt.Sprintf("%s/auth/userinfo", u.baseURL)
	params := map[string]string{
		"token": token,
	}
	return u.makeRequest(endpoint, params)
}

// ValidateUser validates a user token (alias for GetUserInfo)
func (u *UserServiceImpl) ValidateUser(token string) (*models.AuthResponse, error) {
	return u.GetUserInfo(token)
}

// Logout logs out a user
func (u *UserServiceImpl) Logout(token string) error {
	// Since this is a server session system, we log it
	// In a more complex system, you might want to blacklist the token
	fmt.Printf("User logged out with token: %s\n", token)
	return nil
}

// makeRequest is a helper method for user operations
func (u *UserServiceImpl) makeRequest(endpoint string, params map[string]string) (*models.AuthResponse, error) {
	req, err := u.builder.BuildPOSTRequest(endpoint, params)
	if err != nil {
		return nil, err
	}

	_, bodyBytes, err := u.http.Do(req)
	if err != nil {
		return nil, err
	}

	authResp, err := u.parser.ParseAuthResponse(bodyBytes)
	if err != nil {
		return nil, err
	}

	return authResp, u.parser.ValidateResponseSuccess(authResp)
}