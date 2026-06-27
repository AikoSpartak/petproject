package repository

import "arch/internal/models"

type UserRepostory interface {
	Create(user *models.User) error
	Read(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}
