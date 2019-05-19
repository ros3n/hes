package server

import (
	"github.com/ros3n/hes/api/server/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter(emailsHandler *handlers.EmailsAPIHandler) http.Handler {
	router := mux.NewRouter()

	addEmailsAPIHandlers(router, emailsHandler)

	return router
}

func addEmailsAPIHandlers(router *mux.Router, handler *handlers.EmailsAPIHandler) {
	emailsRouter := router.PathPrefix("/emails").Subrouter()
	emailsRouter.HandleFunc("/", handler.CreateEmail).Methods(http.MethodPost)
}
