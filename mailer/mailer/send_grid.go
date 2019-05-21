package mailer

import "github.com/ros3n/hes/mailer/models"

type SendGridMailer struct {
}

func (sg *SendGridMailer) Send(email *models.Email) error {
	return nil
}

type SendGridMailerFactory struct{}

func (lmf *SendGridMailerFactory) Provider() EmailProvider {
	return ProviderSendGrid
}

func (lmf *SendGridMailerFactory) NewMailer() Mailer {
	return &SendGridMailer{}
}
