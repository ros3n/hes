package main

import (
	"github.com/ros3n/hes/mailer/mailer"
	"github.com/ros3n/hes/mailer/manager"
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

	manager := manager.NewManager(nil, mailerFactory)
	manager.Start()

	sig := <-interrupt
	switch sig {
	case os.Interrupt:
	case syscall.SIGTERM:
	}
	manager.Stop()
}
