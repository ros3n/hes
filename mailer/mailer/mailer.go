package mailer

import (
	"errors"

	"github.com/ros3n/hes/mailer/models"
)

type EmailProvider int

const (
	ProviderMailGun EmailProvider = iota
	ProviderSendGrid
	ProviderLocal
)

var ErrProviderNotSupported = errors.New("provider not supported")

type Mailer interface {
	Send(email *models.Email) error
}

type MailerFactory interface {
	Provider() EmailProvider
	NewMailer() Mailer
}

// AbstractMailerFactory contains MailerFactories ready to produce configured clients
type AbstractMailerFactory struct {
	factories map[EmailProvider]MailerFactory
}

// NewMailerFactory groups MailerFactories. Factories passed as arguments should be able to produce correctly configured
// mailers.
func NewMailerFactory(factories ...MailerFactory) *AbstractMailerFactory {
	amf := AbstractMailerFactory{
		factories: make(map[EmailProvider]MailerFactory),
	}
	for _, factory := range factories {
		amf.add(factory)
	}
	return &amf
}

func (amf *AbstractMailerFactory) add(factory MailerFactory) {
	amf.factories[factory.Provider()] = factory
}

// NewMailer returns a configured Mailer that is ready to send emails
func (amf *AbstractMailerFactory) NewMailer(provider EmailProvider) (Mailer, error) {
	mailerFactory, ok := amf.factories[provider]
	if !ok {
		return nil, ErrProviderNotSupported
	}
	return mailerFactory.NewMailer(), nil
}
