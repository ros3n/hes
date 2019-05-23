package server

import (
	"github.com/ros3n/hes/api/server/handlers"
	"github.com/ros3n/hes/api/server/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter(authService middleware.Authenticator, emailsHandler *handlers.EmailsAPIHandler) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.HealthCheck).Methods(http.MethodGet)

	addEmailsAPIHandlers(router, emailsHandler, authService)

	return router
}

func addEmailsAPIHandlers(router *mux.Router, handler *handlers.EmailsAPIHandler, authService middleware.Authenticator) {
	emailsRouter := router.PathPrefix("/emails").Subrouter()
	emailsRouter.Use(middleware.AuthenticationMiddleware(authService))
	emailsRouter.HandleFunc("", handler.ListEmails).Methods(http.MethodGet)
	emailsRouter.HandleFunc("", handler.CreateEmail).Methods(http.MethodPost)
	emailsRouter.HandleFunc("/{id}/send", handler.SendEmail).Methods(http.MethodPost)
}
