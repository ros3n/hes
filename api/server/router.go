package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter(emailsHandler *emailsAPIHandler) http.Handler {
	router := mux.NewRouter()

	addEmailsAPIHandlers(router, emailsHandler)

	return router
}

func addEmailsAPIHandlers(router *mux.Router, handler *emailsAPIHandler) {
	emailsRouter := router.PathPrefix("/emails").Subrouter()
	emailsRouter.HandleFunc("/", handler.createEmail).Methods(http.MethodPost)
}
