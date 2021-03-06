package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ros3n/hes/api/messenger"
	"github.com/ros3n/hes/api/repositories"
	"github.com/ros3n/hes/api/server"
	"github.com/ros3n/hes/api/server/middleware"
	"github.com/ros3n/hes/api/services"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("Starting server..")

	// TODO: use viper
	apiSenderAddr := strings.TrimSpace(os.Getenv("API_SENDER_ADDR"))
	apiReceiverAddr := strings.TrimSpace(os.Getenv("API_RECEIVER_ADDR"))
	apiAddr := strings.TrimSpace(os.Getenv("API_ADDR"))
	dbAddr := strings.TrimSpace(os.Getenv("HES_DATABASE_URL"))
	userName := strings.TrimSpace(os.Getenv("HES_USER_NAME"))
	password := strings.TrimSpace(os.Getenv("HES_PASSWORD"))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	repository, err := repositories.NewDBEmailsRepository(dbAddr)
	if err != nil {
		panic(err)
	}

	//repository := repositories.NewSimpleEmailsRepository()

	msgSender := messenger.NewGRPCMessageSender(apiSenderAddr)
	msgReceiver := messenger.NewGRPCMessageReceiver(apiReceiverAddr)
	emailService := services.NewEmailService(repository, msgSender)

	callbackServer := server.NewCallbackServer(msgReceiver, emailService)

	authService := middleware.NewBasicAuthenticator(userName, password, "1")

	svr := server.NewServer(apiAddr, emailService, authService)
	go func() { log.Fatal(svr.ListenAndServe()) }()
	err = callbackServer.Start()
	if err != nil {
		panic(err)
	}

	log.Println("Server started.")
	log.Println("Waiting for connections..")

	sig := <-interrupt
	switch sig {
	case os.Interrupt:
	case syscall.SIGTERM:
	}

	log.Print("Shutting down..")
	callbackServer.Stop()
	err = svr.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
	log.Println("Done.")
}
