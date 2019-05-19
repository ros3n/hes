package handlers

import (
	"github.com/ros3n/hes/api/server/middleware"
	"net/http"
)

// UserID extracts a user's id from context. If for whatever reason the user id was empty, the app cannot proceed.
func userID(req *http.Request) string {
	rawUserID := req.Context().Value(middleware.UserIDContextKey)
	if rawUserID == nil {
		panic("MISSING USER ID")
	}

	userID := rawUserID.(string)
	if userID == "" {
		panic("EMPTY USER ID")
	}

	return userID
}
