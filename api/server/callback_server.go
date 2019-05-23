package server

import (
	"github.com/ros3n/hes/api/messenger"
	"github.com/ros3n/hes/api/models"
	"github.com/ros3n/hes/api/services"
	"log"
)

type CallbackServer struct {
	emailService    *services.EmailService
	messageReceiver messenger.MessageReceiver
	sendStatusChan  chan *models.SendStatus
	stopChan        chan chan struct{}
}

func NewCallbackServer(msgReceiver messenger.MessageReceiver, emailService *services.EmailService) *CallbackServer {
	sendStatusChan := make(chan *models.SendStatus)
	stopChan := make(chan chan struct{}, 1)
	return &CallbackServer{
		messageReceiver: msgReceiver, emailService: emailService,
		sendStatusChan: sendStatusChan, stopChan: stopChan,
	}
}

func (m *CallbackServer) Start() error {
	log.Println("Starting CallbackServer..")
	m.messageReceiver.Start(m.sendStatusChan)
	go func() {
		for {
			select {
			case status := <-m.sendStatusChan:
				m.emailService.UpdateStatus(status)
			case callback := <-m.stopChan:
				log.Println("CallbackServer is shutting down..")
				m.messageReceiver.Stop()
				log.Println("Done.")
				callback <- struct{}{}
			}
		}
	}()
	log.Println("CallbackServer started.")
	return nil
}

func (m *CallbackServer) Stop() {
	callback := make(chan struct{}, 1)
	m.stopChan <- callback
	<-callback
}
