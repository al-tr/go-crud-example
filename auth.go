package main

import (
	"errors"
	"net/http"
	"strings"
)

// Auth userinfo model
type UserInfo struct {
	Email *string
}

const (
	authorizationHeader = "Authorization"
)

func authenticateRequest(r *http.Request) (*UserInfo, error) {
	headerValue := r.Header.Get(authorizationHeader)
	if stringNilOrEmpty(&headerValue) || !strings.HasPrefix(headerValue, "Bearer ") {
		return nil, errors.New("user is not authenticated")
	}
	emailFromHeader := headerValue[len("Bearer "):]
	user := UserInfo{Email: &emailFromHeader}
	return &user, nil
}
