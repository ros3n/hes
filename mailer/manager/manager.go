package manager

import (
	"github.com/ros3n/hes/mailer/models"
)

type SendStatus struct {
	EmailID int
	Success bool
}

type Manager struct {
	//emailStore   *EmailStore
	callbackChan chan SendStatus
}

func (m *Manager) SendEmail(email *models.Email) error {
	//newEmail, err := m.newEmail(email)
	//if err != nil {
	//	log.Printf("Failed to check presence of email %d", email.ID)
	//	return err
	//}
	//if !newEmail {
	//	log.Printf("Email %d already scheduled for send\n", email.ID)
	//	return err
	//}
	//if err = m.emailStore.Save(email); err != nil {
	//	log.Printf("Failed to store email %d", email.ID)
	//	return err
	//}
	//if err = m.scheduleSend(email); err != nil {
	//
	//}
	return nil
}
