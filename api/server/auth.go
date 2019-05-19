package server

import "errors"

var ErrUserUnauthorized = errors.New("unauthorized")

type Authenticator interface {
	Authenticate(apiKey string) (int, error)
}

type BasicAuthenticator struct {
	apiKey string
	userID int
}

func NewBasicAuthenticator(apiKey string, userID int) *BasicAuthenticator {
	return &BasicAuthenticator{apiKey: apiKey, userID: userID}
}

func (ba *BasicAuthenticator) Authenticate(apiKey string) (int, error) {
	if apiKey == ba.apiKey {
		return ba.userID, nil
	}
	return 0, ErrUserUnauthorized
}