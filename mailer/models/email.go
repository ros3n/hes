package models

// Email struct contains the information necessary to send an email using and external provider.
type Email struct {
	ID         int      `json:"id"`
	Sender     string   `json:"sender"`
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Message    string   `json:"message"`
}

type SendStatus struct {
	EmailID int
	Success bool
}
