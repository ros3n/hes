package models

type EmailStatus string

const (
	EmailCreated EmailStatus = "created"
	EmailQueued              = "queued"
	EmailSent                = "failed"
	EmailBounced             = "bounced"
)

// Email struct contains the information necessary to send an email using and external provider and an information
// about the current status in the pipeline.
type Email struct {
	ID         int         `json:"id"`
	Sender     string      `json:"sender"`
	Recipients []string    `json:"recipients"`
	Subject    string      `json:"subject"`
	Message    string      `json:"message"`
	Status     EmailStatus `json:"status"`
}

// EmailChangeSet is a utility struct that can be used to operate on raw Email data send over HTTP by a client.
type EmailChangeSet struct {
	Sender     *string  `json:"sender"`
	Recipients []string `json:"recipients"`
	Subject    *string  `json:"subject"`
	Message    *string  `json:"message"`
}

func (ecs *EmailChangeSet) ApplyChanges(email *Email) {
	if ecs.Sender != nil {
		email.Sender = *ecs.Sender
	}
	email.Recipients = ecs.Recipients
	if ecs.Subject != nil {
		email.Subject = *ecs.Subject
	}
	if ecs.Message != nil {
		email.Message = *ecs.Message
	}
}
