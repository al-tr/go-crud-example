package auth

import (
	"crud/util"
	"errors"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

func AuthenticateRequest(r *http.Request) (*User, error) {
	headerValue := r.Header.Get(authorizationHeader)
	if util.StringNilOrEmpty(&headerValue) || !strings.HasPrefix(headerValue, "Bearer ") {
		return nil, errors.New("user is not authenticated")
	}
	emailFromHeader := headerValue[len("Bearer "):]
	user := User{Email: &emailFromHeader}
	return &user, nil
}
