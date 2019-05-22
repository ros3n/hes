package manager

import (
	"github.com/ros3n/hes/mailer/mailer"
	"github.com/ros3n/hes/mailer/models"
	"github.com/ros3n/hes/mailer/worker"
	"log"
	"sync"
)

type Manager struct {
	newEmailsChan chan *models.Email     // messenger passes new send requests via this channel
	callbackChan  chan models.SendStatus // workers pass send reports to messenger via this channel
	messenger     Messenger
	factory       *mailer.AbstractMailerFactory
	workerWg      sync.WaitGroup
	stopChan      chan chan struct{}
}

func NewManager(messenger Messenger, factory *mailer.AbstractMailerFactory) *Manager {
	newEmailsChan := make(chan *models.Email)
	callbackChan := make(chan models.SendStatus)
	stopChan := make(chan chan struct{}, 1)
	return &Manager{
		messenger: messenger, factory: factory, newEmailsChan: newEmailsChan, callbackChan: callbackChan,
		workerWg: sync.WaitGroup{}, stopChan: stopChan,
	}
}

func (m *Manager) Start() {
	log.Println("Starting manager..")
	//go m.messenger.Start(newEmailsChan)
	go func() {
		for {
			select {
			case email := <-m.newEmailsChan:
				m.scheduleSend(email)
			case callback := <-m.stopChan:
				log.Println("Manager is shutting down..")
				//m.messenger.Stop()
				m.workerWg.Wait()
				log.Println("Done.")
				callback <- struct{}{}
			}
		}
	}()
	log.Println("Manager started.")
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
