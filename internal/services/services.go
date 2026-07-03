package services

import (
	"arch/internal/models"
	"context"
)

type UserSrvices interface {
	RegisterUser(ctx context.Context, email, password, phone string) (*models.User, error)
	LoginUser(ctx context.Context, login, password string) (*models.User, error)
	GetProfile(ctx context.Context, id int) (*models.User, error)
	UpdateProfile(ctx context.Context, id int, newName, newNumber, newEmail string) (*models.User, error)
	ChangePassword(ctx context.Context, lastPassword, newPassword string)
	DeleteAccount(ctx context.Context, id int) error
}
