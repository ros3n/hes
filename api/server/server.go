package server

import (
	"github.com/gorilla/handlers"
	apiHandlers "github.com/ros3n/hes/api/server/handlers"
	"github.com/ros3n/hes/api/server/middleware"
	"github.com/ros3n/hes/api/services"
	"net/http"
	"os"
)

func NewServer(addr string, emailService *services.EmailService) *http.Server {
	authService := middleware.NewBasicAuthenticator("hypatos", "secret", "1")
	emailsHandler := apiHandlers.NewEmailsAPIHandler(emailService)
	router := newRouter(authService, emailsHandler)
	loggingRouter := handlers.LoggingHandler(os.Stdout, router)

	return &http.Server{
		Addr:    addr,
		Handler: loggingRouter,
	}
}
