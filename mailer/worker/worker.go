package worker

import (
	"github.com/ros3n/hes/mailer/mailer"
	"github.com/ros3n/hes/mailer/models"
	"gopkg.in/matryer/try.v1"
	"log"
	"sync"
)

type Worker struct {
	mailerFactory *mailer.AbstractMailerFactory
	group         sync.WaitGroup
	callbackCh    chan<- *models.SendStatus
	nextProvider  func() mailer.EmailProvider
}

func NewWorker(group sync.WaitGroup, callbackCh chan<- *models.SendStatus,
	mf *mailer.AbstractMailerFactory) *Worker {

	// next function creates a producer that returns EmailProviders using a round-robin strategy
	next := func(providers []mailer.EmailProvider) func() mailer.EmailProvider {
		current := 0
		return func() mailer.EmailProvider {
			provider := providers[current]
			current = (current + 1) % len(providers)
			return provider
		}
	}
	return &Worker{mailerFactory: mf, group: group, callbackCh: callbackCh, nextProvider: next(mf.Providers())}
}

func (w *Worker) SendAsync(email *models.Email) {
	go w.sendEmail(email)
}

func (w *Worker) sendEmail(email *models.Email) {
	maxRetries := 2
	w.group.Add(1)
	defer w.group.Done()

	err := try.Do(func(attempt int) (retry bool, err error) {
		provider := w.nextProvider()
		m, err := w.mailerFactory.NewMailer(provider)
		if err != nil {
			log.Printf("Failed do send email %d: %v\n", email.ID, err)
			return attempt < maxRetries, err
		}
		err = m.Send(email)
		if err != nil {
			log.Printf("Failed do send email %d: %v\n", email.ID, err)
		}
		return attempt < maxRetries, err
	})

	if err != nil {
		log.Printf("Failed to send email %d. Stopping after %d retries\n", email.ID, maxRetries)
		w.callbackCh <- &models.SendStatus{EmailID: email.ID, Success: false}
		return
	}

	log.Printf("Email %d sent successfully\n", email.ID)
	w.callbackCh <- &models.SendStatus{EmailID: email.ID, Success: true}
}
