package repositories

import "github.com/ros3n/hes/api/models"

type EmailsRepository interface {
	Find(id int) (*models.Email, error)
	Create(*models.Email) (*models.Email, error)
	All() ([]*models.Email, error)
}
