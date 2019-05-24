package mailer

import (
	"fmt"

	"github.com/ros3n/hes/mailer/models"
)

type LocalMailer struct {
}

func (lm *LocalMailer) Send(email *models.Email) error {
	fmt.Println(*email)
	return nil
}

type LocalMailerFactory struct{}

func (lmf *LocalMailerFactory) Provider() EmailProvider {
	return ProviderLocal
}

func (lmf *LocalMailerFactory) NewMailer() Mailer {
	return &LocalMailer{}
}
