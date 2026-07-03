package repository

import (
	"arch/internal/models"
	"context"

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

func (u *userRepo) CreateUser(ctx context.Context, user *models.User) error {

	if err := u.DB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepo) FindUserId(ctx context.Context, email string) (*models.User, error) {

	var user models.User
	if err := u.DB.WithContext(ctx).First(user, email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) FindUserEmail(ctx context.Context, email string) (*models.User, error) {

	var user models.User
	if err := u.DB.WithContext(ctx).Where(user, email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, user *models.User) error {

	if err := u.DB.WithContext(ctx).Model(user).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepo) DeleteUser(ctx context.Context, id int) error {
	if err := u.DB.WithContext(ctx).Delete(id).Error; err != nil {
		return err
	}

	return nil
}
