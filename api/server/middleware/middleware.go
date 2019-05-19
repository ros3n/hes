package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

func AuthenticationMiddleware(authService Authenticator) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			userID, err := authService.Authenticate(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := req.Context()
			ctx = context.WithValue(ctx, UserIDContextKey, userID)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
