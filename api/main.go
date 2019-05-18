package main

import (
	"context"
	"github.com/ros3n/hes/api/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting server..")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	svr := server.NewServer("localhost:8080")
	go func () { log.Fatal(svr.ListenAndServe()) }()

	log.Println("Done.")
	log.Println("Waiting for connections..")

	sig := <-interrupt
	switch sig {
	case os.Interrupt:
	case syscall.SIGTERM:
	}

	log.Print("Shutting down..")
	err := svr.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
	log.Print("Done.")
}
