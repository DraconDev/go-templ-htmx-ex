package middleware

import (
	"context"
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
)

// UserContextKey is the key used to store user info in request context
type UserContextKey string

const userContextKey UserContextKey = "user"

// GetUserFromContext gets user info from request context
func GetUserFromContext(r *http.Request) layouts.UserInfo {
	userInfo, ok := r.Context().Value(userContextKey).(layouts.UserInfo)
	if !ok {
		return layouts.UserInfo{LoggedIn: false}
	}
	return userInfo
}

// AddUserToContext adds user info to the request context
func AddUserToContext(r *http.Request, userInfo layouts.UserInfo) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, userInfo)
	return r.WithContext(ctx)
}