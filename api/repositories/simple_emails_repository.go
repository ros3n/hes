package repositories

import (
	"github.com/ros3n/hes/api/models"
	"sync"
)

type SimpleEmailsRepository struct {
	emails          map[int]*models.Email
	autoincrementId int
	mtx             sync.Mutex
}

func NewSimpleEmailsRepository() *SimpleEmailsRepository {
	emails := make(map[int]*models.Email)
	return &SimpleEmailsRepository{emails: emails, autoincrementId: 1}
}

func (ser *SimpleEmailsRepository) Find(id int) (*models.Email, error) {
	ser.mtx.Lock()
	defer ser.mtx.Unlock()

	fetched := ser.emails[id]
	if fetched == nil {
		return nil, nil
	}
	return dup(fetched), nil
}

func (ser *SimpleEmailsRepository) Create(email *models.Email) (*models.Email, error) {
	ser.mtx.Lock()
	defer ser.mtx.Unlock()

	email.ID = ser.autoincrementId
	ser.autoincrementId++
	ser.emails[email.ID] = dup(email)

	return email, nil
}

func (ser *SimpleEmailsRepository) All() ([]*models.Email, error) {
	ser.mtx.Lock()
	defer ser.mtx.Unlock()

	result := make([]*models.Email, 0, len(ser.emails))
	for _, email := range ser.emails {
		result = append(result, dup(email))
	}

	return result, nil
}

func dup(email *models.Email) *models.Email {
	recipients := make([]string, len(email.Recipients))
	copy(recipients, email.Recipients)
	return &models.Email{
		ID:         email.ID,
		Sender:     email.Sender,
		Recipients: recipients,
		Subject:    email.Subject,
		Message:    email.Message,
		Status:     email.Status,
	}
}
