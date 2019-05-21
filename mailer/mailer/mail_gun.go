package mailer

import "github.com/ros3n/hes/mailer/models"

type MailGunMailer struct {
}

func (mg *MailGunMailer) Send(email *models.Email) error {
	return nil
}

type MailGunMailerFactory struct{}

func (lmf *MailGunMailerFactory) Provider() EmailProvider {
	return ProviderMailGun
}

func (lmf *MailGunMailerFactory) NewMailer() Mailer {
	return &MailGunMailer{}
}
