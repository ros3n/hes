package services

import (
	"github.com/ros3n/hes/api/models"
	"github.com/ros3n/hes/api/repositories"
)

type EmailService struct {
	repository repositories.EmailsRepository
}

func NewEmailService(repository repositories.EmailsRepository) *EmailService {
	return &EmailService{repository: repository}
}

func (es *EmailService) UpdateStatus(status *models.SendStatus) error {
	email, err := es.repository.FindByID(status.ID)
	if err != nil {
		return err
	}
	if status.Success {
		email.Status = models.EmailSent
	} else {
		email.Status = models.EmailFailed
	}

	email, err = es.repository.Update(email.UserID, email)
	if err != nil {
		return err
	}
	return nil
}
