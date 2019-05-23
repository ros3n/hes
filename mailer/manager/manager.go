package manager

import (
	"context"
	"github.com/ros3n/hes/mailer/mailer"
	"github.com/ros3n/hes/mailer/messenger"
	"github.com/ros3n/hes/mailer/models"
	"github.com/ros3n/hes/mailer/worker"
	"log"
	"sync"
)

type Manager struct {
	newEmailsChan   chan *models.Email      // messageReceiver passes new send requests via this channel
	callbackChan    chan *models.SendStatus // workers pass send reports to messageReceiver via this channel
	messageReceiver messenger.MessageReceiver
	messageSender   messenger.MessageSender
	factory         *mailer.AbstractMailerFactory
	workerWg        sync.WaitGroup
	stopChan        chan chan struct{}
}

func NewManager(msgReceiver messenger.MessageReceiver, msgSender messenger.MessageSender,
	factory *mailer.AbstractMailerFactory) *Manager {

	newEmailsChan := make(chan *models.Email)
	callbackChan := make(chan *models.SendStatus)
	stopChan := make(chan chan struct{}, 1)
	return &Manager{
		messageReceiver: msgReceiver, messageSender: msgSender, factory: factory, newEmailsChan: newEmailsChan,
		callbackChan: callbackChan, workerWg: sync.WaitGroup{}, stopChan: stopChan,
	}
}

func (m *Manager) Start() error {
	log.Println("Starting manager..")
	err := m.messageReceiver.Start(m.newEmailsChan)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case email := <-m.newEmailsChan:
				m.scheduleSend(email)
			case status := <-m.callbackChan:
				m.notifyAPI(status)
			case callback := <-m.stopChan:
				log.Println("Manager is shutting down..")
				m.messageReceiver.Stop()
				m.workerWg.Wait()
				log.Println("Done.")
				callback <- struct{}{}
			}
		}
	}()

	log.Println("Manager started.")
	return nil
}

func (m *Manager) Stop() {
	callback := make(chan struct{}, 1)
	m.stopChan <- callback
	<-callback
}

func (m *Manager) scheduleSend(email *models.Email) {
	wkr := worker.NewWorker(m.workerWg, m.callbackChan, m.factory)
	wkr.SendAsync(email)
}

func (m *Manager) notifyAPI(status *models.SendStatus) {
	// TODO: notify asynchronously
	err := m.messageSender.SendStatus(context.Background(), status)
	if err != nil {
		log.Println(err)
	}
}
