package repositories

import (
	"errors"
	"github.com/ros3n/hes/api/models"
)

var (
	ErrFetchFailed   = errors.New("failed to fetch data")
	ErrCreateFailed  = errors.New("failed to insert data")
	ErrEmailNotFound = errors.New("email not found")
)

type EmailsRepository interface {
	Find(userID string, id int) (*models.Email, error)
	Create(userID string, email *models.Email) (*models.Email, error)
	Update(userID string, email *models.Email) (*models.Email, error)
	All(userID string) (emails []*models.Email, err error)
}
