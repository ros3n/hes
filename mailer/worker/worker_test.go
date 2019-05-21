package worker

import (
	"errors"
	"github.com/ros3n/hes/mailer/mailer"
	"github.com/ros3n/hes/mailer/manager"
	"github.com/ros3n/hes/mailer/models"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

const emailID = 1

type WorkerTestSuite struct {
	suite.Suite
	wg         sync.WaitGroup
	callbackCh chan manager.SendStatus
	email      *models.Email
}

func (ws *WorkerTestSuite) SetupSuite() {
	log.SetOutput(ioutil.Discard)
}

func (ws *WorkerTestSuite) AfterSuite() {
	log.SetOutput(os.Stdout)
}

func (ws *WorkerTestSuite) SetupTest() {
	ws.wg = sync.WaitGroup{}
	ws.callbackCh = make(chan manager.SendStatus)
	ws.email = &models.Email{ID: emailID}
}

func (ws *WorkerTestSuite) TestSuccessfulSend() {
	mailerFactory := ws.setupMailers(false, func(_ *bool) {})
	worker := ws.setupWorker(mailerFactory)

	worker.SendAsync(ws.email)

	select {
	case status := <-ws.callbackCh:
		ws.Equal(ws.email.ID, status.EmailID)
		ws.True(status.Success)
	case <-time.After(time.Second):
		ws.Fail("Worker freezed.")
	}
}

func (ws *WorkerTestSuite) TestSuccessfulSendAfterRetry() {
	mailerFactory := ws.setupMailers(true, func(faulty *bool) {
		if *faulty {
			*faulty = false
		}
	})
	worker := ws.setupWorker(mailerFactory)

	worker.SendAsync(ws.email)

	select {
	case status := <-ws.callbackCh:
		ws.Equal(ws.email.ID, status.EmailID)
		ws.True(status.Success)
	case <-time.After(time.Second):
		ws.Fail("Worker freezed.")
	}
}

func (ws *WorkerTestSuite) TestFailAfterRetry() {
	mailerFactory := ws.setupMailers(true, func(faulty *bool) {})
	worker := ws.setupWorker(mailerFactory)

	worker.SendAsync(ws.email)

	select {
	case status := <-ws.callbackCh:
		ws.Equal(ws.email.ID, status.EmailID)
		ws.False(status.Success)
	case <-time.After(time.Second):
		ws.Fail("Worker freezed.")
	}
}

func TestRunWorkerSuite(t *testing.T) {
	suite.Run(t, new(WorkerTestSuite))
}

func (ws *WorkerTestSuite) setupWorker(mailerFactory *mailer.AbstractMailerFactory) *Worker {
	return NewWorker(ws.wg, ws.callbackCh, mailerFactory, mailer.ProviderLocal)
}

func (ws *WorkerTestSuite) setupMailers(faulty bool, updateFaultyStrategy func(*bool)) *mailer.AbstractMailerFactory {
	tmf := &TestMailerFactory{faulty: faulty, updateFaulty: updateFaultyStrategy}
	amf := mailer.NewMailerFactory(tmf)
	return amf
}

type TestMailer struct {
	faulty bool
}

func (tm *TestMailer) Send(email *models.Email) error {
	if tm.faulty {
		return errors.New("failed to send the email")
	}
	return nil
}

type TestMailerFactory struct {
	faulty       bool
	updateFaulty func(*bool)
}

func (tmf *TestMailerFactory) Provider() mailer.EmailProvider {
	return mailer.ProviderLocal
}

func (tmf *TestMailerFactory) NewMailer() mailer.Mailer {
	defer tmf.updateFaulty(&tmf.faulty)
	return &TestMailer{faulty: tmf.faulty}
}
