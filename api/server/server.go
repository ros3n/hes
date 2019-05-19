package server

import (
	"github.com/gorilla/handlers"
	"github.com/ros3n/hes/api/repositories"
	apiHandlers "github.com/ros3n/hes/api/server/handlers"
	"github.com/ros3n/hes/api/server/middleware"
	"net/http"
	"os"
)

func NewServer(addr string) *http.Server {
	authService := middleware.NewBasicAuthenticator("hypatos", "secret", "1")
	repo := repositories.NewSimpleEmailsRepository()
	emailsHandler := apiHandlers.NewEmailsAPIHandler(repo)
	router := newRouter(authService, emailsHandler)
	loggingRouter := handlers.LoggingHandler(os.Stdout, router)

	return &http.Server{
		Addr:    addr,
		Handler: loggingRouter,
	}
}
