package mailer

import (
	"fmt"
	"net/http"

	"github.com/ros3n/hes/mailer/models"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	apiKey string
}

func (sg *SendGridMailer) Send(email *models.Email) error {
	message := sg.composeEmail(email)
	client := sendgrid.NewSendClient(sg.apiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("Got a non-202 response: %v", response)
		return err
	}

	return nil
}

func (sg *SendGridMailer) composeEmail(email *models.Email) *mail.SGMailV3 {
	from := mail.NewEmail("", email.Sender)
	subject := email.Subject
	to := mail.NewEmail("", email.Recipients[0])
	plainTextContent := email.Message
	htmlContent := email.Message

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	addAdditionalRecipients(message, email.Recipients[1:])

	return message
}

func addAdditionalRecipients(message *mail.SGMailV3, recipients []string) {
	if len(recipients) == 0 {
		return
	}
	personalization := mail.NewPersonalization()
	for _, additionalRecipient := range recipients {
		personalization.AddTos(mail.NewEmail("", additionalRecipient))
	}
	message.AddPersonalizations(personalization)
	return
}

type SendGridMailerFactory struct {
	apiKey string
}

func NewSendGridMailerFactory(apiKey string) *SendGridMailerFactory {
	return &SendGridMailerFactory{apiKey: apiKey}
}

func (lmf *SendGridMailerFactory) Provider() EmailProvider {
	return ProviderSendGrid
}

func (lmf *SendGridMailerFactory) NewMailer() Mailer {
	return &SendGridMailer{apiKey: lmf.apiKey}
}
