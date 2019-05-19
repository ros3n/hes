package server

import (
	"github.com/ros3n/hes/api/server/handlers"
	"github.com/ros3n/hes/api/server/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter(authService middleware.Authenticator, emailsHandler *handlers.EmailsAPIHandler) http.Handler {
	router := mux.NewRouter()
	router.Use(middleware.AuthenticationMiddleware(authService))

	addEmailsAPIHandlers(router, emailsHandler)

	return router
}

func addEmailsAPIHandlers(router *mux.Router, handler *handlers.EmailsAPIHandler) {
	emailsRouter := router.PathPrefix("/emails").Subrouter()
	emailsRouter.HandleFunc("/", handler.ListEmails).Methods(http.MethodGet)
	emailsRouter.HandleFunc("/", handler.CreateEmail).Methods(http.MethodPost)
}
