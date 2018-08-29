package main

import (
	"errors"
	"net/http"
	"strings"
)

type User struct {
	Email *string
}

const (
	authorizationHeader = "Authorization"
)

func authenticateRequest(r *http.Request) (*User, error) {
	headerValue := r.Header.Get(authorizationHeader)
	if stringNilOrEmpty(&headerValue) || !strings.HasPrefix(headerValue, "Bearer ") {
		return nil, errors.New("user is not authenticated")
	}
	emailFromHeader := headerValue[len("Bearer "):]
	user := User{Email: &emailFromHeader}
	return &user, nil
}
