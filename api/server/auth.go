package server

import (
	"errors"
	"net/http"
)

var ErrUserUnauthorized = errors.New("unauthorized")

type Authenticator interface {
	Authenticate(req *http.Request) (string, error)
}

type BasicAuthenticator struct {
	userName string
	password string
	userID string
}

func NewBasicAuthenticator(userName, password, userID string) *BasicAuthenticator {
	return &BasicAuthenticator{userName: userName, password: password, userID: userID}
}

func (ba *BasicAuthenticator) Authenticate(req *http.Request) (string, error) {
	username, password, ok := req.BasicAuth()
	if !ok {
		return "", ErrUserUnauthorized
	}

	if username != ba.userName || password != ba.password {
		return "", ErrUserUnauthorized
	}

	return ba.userID, nil
}