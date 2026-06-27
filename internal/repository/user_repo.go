package repository

import (
	"arch/internal/models"

	"gorm.io/gorm"
)

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepo {
	return &userRepo{
		DB: db,
	}
}

func (u *userRepo) Create(user *models.User) error {
	return u.DB.Create(user).Error
}

func (u *userRepo) Read(id int) (*models.User, error) {

	var user models.User
	err := u.DB.First(&user, id).Error

	return &user, err
}

func (u *userRepo) Update(user *models.User) error {
	return u.DB.Save(user).Error
}

func (u *userRepo) Delete(id int) error {

	return u.DB.Delete(&models.User{}, id).Error
}
