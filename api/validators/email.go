package validators

import (
	"fmt"
	"regexp"

	"github.com/ros3n/hes/api/models"
)

var emailRegEx = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

const (
	senderKey     = "sender"
	recipientsKey = "recipients"
	subjectKey    = "subject"
	messageKey    = "message"
)

// EmailValidator makes sure that a payload deserialized to an EmailChangeSet doesn't contain blank fields
// and that the data is formatted correctly.
type EmailValidator struct {
	*BaseValidator
	changeSet *models.EmailChangeSet
}

func NewEmailValidator(changeSet *models.EmailChangeSet) *EmailValidator {
	return &EmailValidator{&BaseValidator{}, changeSet}
}

func (ev *EmailValidator) Validate() {
	ev.validateSender()
	ev.validateRecipients()
	ev.validateSubject()
	ev.validateMessage()
}

func (ev *EmailValidator) validateSender() {
	sender := ev.changeSet.Sender
	if !validateStringPresence(sender) {
		ev.addError(senderKey, cannotBeBlankError)
		return
	}
	if !validateEmailAddress(*sender) {
		ev.addError(senderKey, invalidEmailError(*sender))
	}
}

func (ev *EmailValidator) validateRecipients() {
	recipients := ev.changeSet.Recipients
	if recipients == nil || len(recipients) == 0 {
		ev.addError(recipientsKey, cannotBeBlankError)
		return
	}
	for _, recipient := range recipients {
		if !validateEmailAddress(recipient) {
			ev.addError(recipientsKey, invalidEmailError(recipient))
		}
	}
}

func (ev *EmailValidator) validateSubject() {
	subject := ev.changeSet.Subject
	if !validateStringPresence(subject) {
		ev.addError(subjectKey, cannotBeBlankError)
		return
	}
}

func (ev *EmailValidator) validateMessage() {
	message := ev.changeSet.Message
	if !validateStringPresence(message) {
		ev.addError(messageKey, cannotBeBlankError)
	}
}

func validateStringPresence(str *string) bool {
	return str != nil && len(*str) > 0
}

func validateEmailAddress(email string) bool {
	if len(email) >= 255 || !emailRegEx.MatchString(email) {
		return false
	}
	return true
}

func invalidEmailError(email string) string {
	return fmt.Sprintf("%s is not a valid email address", email)
}
