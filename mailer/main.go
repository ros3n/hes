package main

import (
	"github.com/ros3n/hes/mailer/mailer"
	"github.com/ros3n/hes/mailer/manager"
	"github.com/ros3n/hes/mailer/messenger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.Lshortfile)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	sendGridFactory := mailer.NewSendGridMailerFactory(os.Getenv("SEND_GRID_API_KEY"))
	mailGunFactory := mailer.NewMailGunMailerFactory(
		os.Getenv("MAIL_GUN_DOMAIN"), os.Getenv("MAIL_GUN_API_KEY"),
	)
	mailerFactory := mailer.NewMailerFactory(mailGunFactory, sendGridFactory)
	msgReceiver := messenger.NewGRPCMessageReceiver(":8888")

	manager := manager.NewManager(msgReceiver, mailerFactory)
	err := manager.Start()
	if err != nil {
		panic(err)
	}

	sig := <-interrupt
	switch sig {
	case os.Interrupt:
	case syscall.SIGTERM:
	}
	manager.Stop()
}
