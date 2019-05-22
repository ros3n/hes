package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/ros3n/hes/api/server/middleware"
	"net/http"
	"strconv"
)

var (
	ErrNotFound    = errors.New("requested resource was not found")
	ErrServerError = errors.New("internal server error")
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

func emailID(req *http.Request) int64 {
	strId := mux.Vars(req)["id"]
	id, _ := strconv.Atoi(strId)
	return int64(id)
}
