package server

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

func authenticationMiddleware(authService Authenticator) func(handler http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			apiKey := extractApiKey(req)
			userID, err := authService.Authenticate(apiKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := req.Context()
			ctx = context.WithValue(ctx, UserIDContextKey, userID)
			next.ServeHTTP(w, req)
		})
	}
}

func extractApiKey(req *http.Request) string {
	return req.Header.Get("X-API-Key")
}