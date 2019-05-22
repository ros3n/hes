package repositories

import (
	"github.com/ros3n/hes/api/models"
	"sync"
)

type SimpleEmailsRepository struct {
	emails          map[string]map[int64]*models.Email
	autoincrementId int64
	mtx             sync.Mutex
}

func NewSimpleEmailsRepository() *SimpleEmailsRepository {
	emails := make(map[string]map[int64]*models.Email)
	return &SimpleEmailsRepository{emails: emails, autoincrementId: 1}
}

func (ser *SimpleEmailsRepository) Find(userID string, id int64) (*models.Email, error) {
	ser.mtx.Lock()
	defer ser.mtx.Unlock()

	if ser.emails[userID] == nil {
		return nil, nil
	}
	fetched := ser.emails[userID][id]
	if fetched == nil {
		return nil, nil
	}
	return dup(fetched), nil
}

func (ser *SimpleEmailsRepository) Create(userID string, email *models.Email) (*models.Email, error) {
	ser.mtx.Lock()
	defer ser.mtx.Unlock()

	email.ID = ser.autoincrementId
	ser.autoincrementId++
	email.UserID = userID

	if ser.emails[userID] == nil {
		ser.emails[userID] = make(map[int64]*models.Email)
	}
	ser.emails[userID][email.ID] = dup(email)

	return email, nil
}

func (ser *SimpleEmailsRepository) Update(userID string, email *models.Email) (*models.Email, error) {
	ser.mtx.Lock()
	defer ser.mtx.Unlock()

	if ser.emails[userID] == nil {
		return nil, ErrEmailNotFound
	}
	if ser.emails[userID][email.ID] == nil {
		return nil, ErrEmailNotFound
	}
	ser.emails[userID][email.ID] = dup(email)

	return email, nil
}

func (ser *SimpleEmailsRepository) All(userID string) ([]*models.Email, error) {
	ser.mtx.Lock()
	defer ser.mtx.Unlock()

	emails := ser.emails[userID]
	if emails == nil {
		return []*models.Email{}, nil
	}

	result := make([]*models.Email, 0, len(emails))
	for _, email := range emails {
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
