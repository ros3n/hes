package repositories

import (
	"errors"
	"github.com/ros3n/hes/api/models"
)

var (
	ErrFetchFailed  = errors.New("failed to fetch data")
	ErrCreateFailed = errors.New("failed to insert data")
)

type EmailsRepository interface {
	Find(id int) (*models.Email, error)
	Create(*models.Email) (*models.Email, error)
	All() ([]*models.Email, error)
}
