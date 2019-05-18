package server

import (
	"github.com/gorilla/handlers"
	"net/http"
	"os"
)

func NewServer(addr string) *http.Server {
	emailsHandler := newEmailsAPIHandler()
	router := newRouter(emailsHandler)
	loggingRouter := handlers.LoggingHandler(os.Stdout, router)

	return &http.Server{
		Addr:    addr,
		Handler: loggingRouter,
	}
}
