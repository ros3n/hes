package main

import (
	"fmt"
	"github.com/ros3n/hes/mailer/mailer"
	"github.com/ros3n/hes/mailer/manager"
	"github.com/ros3n/hes/mailer/models"
	"github.com/ros3n/hes/mailer/worker"
	"os"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	callbackCh := make(chan manager.SendStatus)

	sendGridFactory := mailer.NewSendGridMailerFactory(os.Getenv("SEND_GRID_API_KEY"))
	mailGunFactory := mailer.NewMailGunMailerFactory(
		os.Getenv("MAIL_GUN_DOMAIN"), os.Getenv("MAIL_GUN_API_KEY"),
	)
	amfMailGun := mailer.NewMailerFactory(mailGunFactory)
	amfSendGrid := mailer.NewMailerFactory(sendGridFactory)

	workerSendGrid := worker.NewWorker(wg, callbackCh, amfSendGrid, mailer.ProviderSendGrid)
	workerMailGun := worker.NewWorker(wg, callbackCh, amfMailGun, mailer.ProviderMailGun)

	var worker *worker.Worker
	var sender string

	t := os.Args[1]
	if t == "sg" {
		worker = workerSendGrid
		sender = "send.grid.test@sendgird.com"
	} else {
		worker = workerMailGun
		sender = fmt.Sprintf("mail.gun@%s", os.Getenv("MAIL_GUN_DOMAIN"))
	}

	email := &models.Email{
		Sender:     sender,
		Recipients: os.Args[2:],
		Subject:    "test",
		Message:    "test message",
	}

	worker.SendAsync(email)

	select {
	case status := <-callbackCh:
		fmt.Println(status)
	case <-time.After(time.Second * 5):
		fmt.Println("TIMEOUT")
	}
}
