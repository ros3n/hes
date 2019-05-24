package server

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	apiHandlers "github.com/ros3n/hes/api/server/handlers"
	"github.com/ros3n/hes/api/server/middleware"
	"github.com/ros3n/hes/api/services"
)

func NewServer(addr string, emailService *services.EmailService, authService middleware.Authenticator) *http.Server {
	emailsHandler := apiHandlers.NewEmailsAPIHandler(emailService)
	router := newRouter(authService, emailsHandler)
	loggingRouter := handlers.LoggingHandler(os.Stdout, router)

	return &http.Server{
		Addr:    addr,
		Handler: loggingRouter,
	}
}
