package main

import (
	"github.com/ros3n/hes/mailer/mailer"
	"github.com/ros3n/hes/mailer/manager"
	"github.com/ros3n/hes/mailer/messenger"
	"github.com/ros3n/hes/mailer/server"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	log.SetFlags(log.Lshortfile)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	sendGridKey := strings.TrimSpace(os.Getenv("SEND_GRID_API_KEY"))
	mailGunDomain := strings.TrimSpace(os.Getenv("MAIL_GUN_DOMAIN"))
	mailGunKey := strings.TrimSpace(os.Getenv("MAIL_GUN_API_KEY"))
	receiverAddr := strings.TrimSpace(os.Getenv("MAILER_RECEIVER_ADDR"))
	senderAddr := strings.TrimSpace(os.Getenv("MAILER_SENDER_ADDR"))
	healthCheckAddr := strings.TrimSpace(os.Getenv("MAILER_HEALTH_CHECK_ADDR"))

	sendGridFactory := mailer.NewSendGridMailerFactory(sendGridKey)
	mailGunFactory := mailer.NewMailGunMailerFactory(mailGunDomain, mailGunKey)
	mailerFactory := mailer.NewMailerFactory(mailGunFactory, sendGridFactory)
	msgReceiver := messenger.NewGRPCMessageReceiver(receiverAddr)
	msgSender := messenger.NewGRPCMessageSender(senderAddr)

	manager := manager.NewManager(msgReceiver, msgSender, mailerFactory)
	err := manager.Start()
	if err != nil {
		panic(err)
	}

	go server.StartHealthCheckServer(healthCheckAddr)

	sig := <-interrupt
	switch sig {
	case os.Interrupt:
	case syscall.SIGTERM:
	}
	manager.Stop()
}
