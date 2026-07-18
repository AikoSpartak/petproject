package services

import (
	"arch/internal/models"
	"arch/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorSenEmail   = errors.New("email format error")
	ErrorsLogin     = errors.New("login format error")
	ErrorsPassword  = errors.New("password format error")
	ErrorsPhone     = errors.New("phone format error")
	ErrUserNotFound = errors.New("user not found")
	ErrUserUuid     = errors.New("uuid not found")
)

type UserSrvices interface {
	RegisterUser(ctx context.Context, login, password, phone, datebirthday, passport, country string) (*models.User, error)
	LoginUser(ctx context.Context, login, password string) (*models.User, error)
	GetProfile(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, name, phone string) (*models.User, error)
	ChangeEmail(ctx context.Context, id uuid.UUID, newEmail string) (*models.User, error)
	ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

type UserService struct {
	UserRep repository.UserRepository
}

func (u *UserService) RegisterUser(ctx context.Context, login, password, phone, datebirthday, passport, country string) (*models.User, error) {

	err := emailValidate(login)
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

	findEmail, err := u.UserRep.GetByEmail(ctx, login)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("the user was not found")
	}

	if findEmail != nil {
		return nil, fmt.Errorf("the user already exists")
	}

	user := &models.User{
		Uuid:        uuid.New(),
		Email:       login,
		Password:    hashPass,
		PhoneNumber: phone,
	}

	if err := u.UserRep.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	return user, nil

}

func (u *UserService) LoginUser(ctx context.Context, login, password string) (*models.User, error) {

	user, err := u.UserRep.GetByEmail(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("invalid username")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

func (u *UserService) GetProfile(ctx context.Context, id uuid.UUID) (*models.User, error) {

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) UpdateProfile(ctx context.Context, id uuid.UUID, name, phone string) (*models.User, error) {

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.FullName = name
	user.PhoneNumber = phone

	u.UserRep.UpdateUser(ctx, user)

	return user, nil
}

func (u *UserService) ChangeEmail(ctx context.Context, id uuid.UUID, newEmail string) (*models.User, error) {

	err := emailValidate(newEmail)
	if err != nil {
		return nil, ErrorSenEmail
	}

	user, err := u.UserRep.GetByEmail(ctx, newEmail)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user != nil && user.Uuid != id {
		return nil, err
	}

	userEmail, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		return nil, err
	}

	userEmail.Email = newEmail

	err = u.UserRep.UpdateUser(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return userEmail, nil

}

func (u *UserService) ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error {

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("password doesn't match")
	}

	hashPass, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashPass

	return nil

}

func (u *UserService) DeleteAccount(ctx context.Context, id uuid.UUID) error {

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		return err
	}

	return u.UserRep.DeleteUser(ctx, user.Uuid)
}
