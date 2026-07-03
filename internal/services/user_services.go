package services

import (
	"arch/internal/models"
	"arch/internal/repository"
	"context"
	"errors"
	"fmt"
)

var (
	ErrorSenEmail   = errors.New("email format error")
	ErrorsLogin     = errors.New("login format error")
	ErrorsPassword  = errors.New("password format error")
	ErrorsPhone     = errors.New("phone format error")
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	userRep repository.UserRepository
}

func NewUserService(userRep repository.UserRepository) *UserService {
	return &UserService{
		userRep: userRep,
	}
}

func (u *UserService) RegisterUser(ctx context.Context, email, password, phone string) (*models.User, error) {

	err := emailValidate(email)
	if err != nil {
		return nil, err
	}

	hashPass, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	err = validateNumber(phone)
	if err != nil {
		return nil, err
	}

	findEmail, err := u.userRep.FindUserEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("the user was not found")
	}

	if findEmail != nil {
		return nil, fmt.Errorf("the user already exists")
	}

	user := &models.User{
		Email:       email,
		Password:    hashPass,
		PhoneNumber: phone,
	}

	if err := u.userRep.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	return user, nil

}
