package main

import (
	"context"
	"github.com/ros3n/hes/api/messenger"
	"github.com/ros3n/hes/api/repositories"
	"github.com/ros3n/hes/api/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("Starting server..")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	repository, err := repositories.NewDBEmailsRepository(os.Getenv("HES_DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	msgSender := messenger.NewGRPCMessageSender("localhost:8888")

	svr := server.NewServer("localhost:8080", repository, msgSender)
	go func() { log.Fatal(svr.ListenAndServe()) }()

	log.Println("Done.")
	log.Println("Waiting for connections..")

	sig := <-interrupt
	switch sig {
	case os.Interrupt:
	case syscall.SIGTERM:
	}

	log.Print("Shutting down..")
	err = svr.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
	log.Print("Done.")
}
