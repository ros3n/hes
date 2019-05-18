package models

type EmailStatus string

const (
	created EmailStatus = "created"
	queued              = "queued"
	sent                = "failed"
	bounced             = "bounced"
)

// Email struct contains the information necessary to send an email using and external provider and an information
// about the current status in the pipeline.
type Email struct {
	Sender     string
	Recipients []string
	Subject    string
	Message    string
	Status     EmailStatus
}

// EmailChangeSet is a utility struct that can be used to operate on raw Email data send over HTTP by a client.
type EmailChangeSet struct {
	Sender     *string  `json:"sender"`
	Recipients []string `json:"recipients"`
	Subject    *string  `json:"subject"`
	Message    *string  `json:"message"`
}
