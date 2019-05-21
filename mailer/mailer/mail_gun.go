package mailer

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v3"
	"github.com/ros3n/hes/mailer/models"
)

type MailGunMailer struct {
	domain string
	apiKey string
}

func (mg *MailGunMailer) Send(email *models.Email) error {
	client := mailgun.NewMailgun(mg.domain, mg.apiKey)

	message := client.NewMessage(email.Sender, email.Subject, email.Message, email.Recipients...)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := client.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

type MailGunMailerFactory struct {
	domain string
	apiKey string
}

func NewMailGunMailerFactory(domain, apiKey string) *MailGunMailerFactory {
	return &MailGunMailerFactory{domain: domain, apiKey: apiKey}
}

func (lmf *MailGunMailerFactory) Provider() EmailProvider {
	return ProviderMailGun
}

func (lmf *MailGunMailerFactory) NewMailer() Mailer {
	return &MailGunMailer{
		domain: lmf.domain,
		apiKey: lmf.apiKey,
	}
}
