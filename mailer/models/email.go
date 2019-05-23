package models

// Email struct contains the information necessary to send an email using and external provider.
type Email struct {
	ID         int64
	Sender     string
	Recipients []string
	Subject    string
	Message    string
}

type SendStatus struct {
	EmailID int64
	Success bool
}
