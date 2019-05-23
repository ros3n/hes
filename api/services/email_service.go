package services

import (
	"context"
	"log"

	"github.com/ros3n/hes/api/messenger"
	"github.com/ros3n/hes/api/models"
	"github.com/ros3n/hes/api/repositories"
	"github.com/ros3n/hes/api/validators"
)

type EmailService struct {
	repository repositories.EmailsRepository
	msgSender  messenger.MessageSender
}

func NewEmailService(repository repositories.EmailsRepository, sender messenger.MessageSender) *EmailService {
	return &EmailService{repository: repository, msgSender: sender}
}

func (es *EmailService) Create(userID string, changeSet *models.EmailChangeSet) (*models.Email, error) {

	validator := validators.NewEmailValidator(changeSet)
	validator.Validate()

	if !validator.Valid() {
		return nil, validator
	}

	email := &models.Email{Status: "created"}
	changeSet.ApplyChanges(email)

	email, err := es.repository.Create(userID, email)
	if err != nil {
		return nil, err
	}

	return email, nil
}

func (es *EmailService) Find(userID string, emailID int64) (*models.Email, error) {
	return es.repository.Find(userID, emailID)
}

func (es *EmailService) All(userID string) ([]*models.Email, error) {
	return es.repository.All(userID)
}

func (es *EmailService) UpdateStatus(status *models.SendStatus) error {
	email, err := es.repository.FindByID(status.ID)
	if err != nil {
		return err
	}
	if status.Success {
		email.Status = models.EmailSent
		log.Printf("Email %d has been sent.", email.ID)
	} else {
		email.Status = models.EmailFailed
		log.Printf("Email %d has failed.", email.ID)
	}

	email, err = es.repository.Update(email.UserID, email)
	if err != nil {
		return err
	}
	return nil
}

func (es *EmailService) ScheduleSend(ctx context.Context, email *models.Email) error {
	log.Printf("Scheduling send for email %d\n", email.ID)
	err := es.msgSender.SendEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return err
	}

	email.Status = models.EmailQueued
	email, err = es.repository.Update(email.UserID, email)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
