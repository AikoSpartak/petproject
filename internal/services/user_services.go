package services

import (
	"arch/internal/models"
	"arch/internal/repository"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=user_services.go -destination=mocks/mock.go

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
	LoginUser(ctx context.Context, login, password string) (string, string, *models.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, name, phone string) (*models.User, error)
	ChangeEmailUser(ctx context.Context, id uuid.UUID, newEmail string) (*models.User, error)
	ChangePasswordUser(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error
	DeleteAccountUser(ctx context.Context, id uuid.UUID) error
}

type UserService struct {
	UserRep repository.UserRepository
}

func (u *UserService) RegisterUser(ctx context.Context, login, password, phone, datebirthday, passport, country string) (*models.User, error) {

	err := emailValidate(login)
	if err != nil {
		slog.Error("mail validation failed",
			"login", login,
			"error", err.Error())
		return nil, err
	}

	hashPass, err := hashPassword(password)
	if err != nil {
		slog.Error("the password is incorrect",
			"password", password,
			"error", err.Error())
		return nil, err
	}

	err = validateNumber(phone)
	if err != nil {
		slog.Error("phone validation failed",
			"phone", phone,
			"error", err.Error())
		return nil, err
	}

	findEmail, err := u.UserRep.GetByEmail(ctx, login)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		slog.Error("couldn't get the user",
			"email", login,
			"error", err.Error())
		return nil, fmt.Errorf("the user was not found")
	}

	if findEmail != nil {
		slog.Warn("user already exists",
			"email", login)
		return nil, fmt.Errorf("the user already exists")
	}

	user := &models.User{
		Uuid:        uuid.New(),
		Email:       login,
		Password:    hashPass,
		PhoneNumber: phone,
	}

	if err := u.UserRep.CreateUser(ctx, user); err != nil {
		slog.Error("failed couldn't add user",
			"email", login,
			"error", err.Error())
		return nil, fmt.Errorf("failed to create user")
	}

	return user, nil

}

func (u *UserService) LoginUser(ctx context.Context, login, password string) (string, string, *models.User, error) {

	user, err := u.UserRep.GetByEmail(ctx, login)
	if err != nil {
		slog.Warn("couldn't get the user",
			"email", login,
			"err", err.Error())
		return "", "", nil, fmt.Errorf("invalid username")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		slog.Error("invalid password",
			"email", login,
			"error", err.Error())
		return "", "", nil, fmt.Errorf("invalid password")
	}

	accessToken, err := GenerateAccessToken(
		user.Uuid.String(),
		user.Role,
	)
	if err != nil {
		slog.Error("failed to generate jwt",
			"email", login,
			"error", err.Error(),
		)
		return "", "", nil, err
	}

	refreshtoken, err := GenerateRefreshToken(
		user.Uuid.String(),
	)
	if err != nil {
		slog.Error("failed to generate jwt",
			"email", login,
			"error", err.Error(),
		)
		return "", "", nil, err
	}

	refreshTokenModel := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.Uuid,
		Token:     refreshtoken,
		ExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}

	if err := u.UserRep.SaveRefreshToken(ctx, refreshTokenModel); err != nil {
		slog.Error("failed to save refresh token",
			"user_id", user.Uuid,
			"error", err.Error())
	}

	return accessToken, refreshtoken, user, nil
}

func (u *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		slog.Warn("couldn't find the user",
			"id", id,
			"error", err.Error())
		return nil, err
	}

	return user, nil
}

func (u *UserService) UpdateUser(ctx context.Context, id uuid.UUID, name, phone string) (*models.User, error) {

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		slog.Error("couldn't find the user",
			"id", id,
			"error", err.Error())
		return nil, ErrUserNotFound
	}

	user.FullName = name
	user.PhoneNumber = phone

	u.UserRep.UpdateUser(ctx, user)

	return user, nil
}

func (u *UserService) ChangeEmailUser(ctx context.Context, id uuid.UUID, newEmail string) (*models.User, error) {

	err := emailValidate(newEmail)
	if err != nil {
		slog.Error("mail validation failed",
			"email", newEmail,
			"error", err.Error())
		return nil, ErrorSenEmail
	}

	user, err := u.UserRep.GetByEmail(ctx, newEmail)
	if err != nil {
		slog.Error("couldn't get the user email",
			"error", err.Error())
		return nil, ErrUserNotFound
	}

	if user != nil && user.Uuid != id {
		return nil, err
	}

	userEmail, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		slog.Error("couldn't find the user id",
			"id", id,
			"error", err.Error())
		return nil, err
	}

	userEmail.Email = newEmail

	err = u.UserRep.UpdateUser(ctx, userEmail)
	if err != nil {
		slog.Error("couldn't update data",
			"error", err.Error(),
			"email", newEmail)
		return nil, err
	}

	return userEmail, nil

}

func (u *UserService) ChangePasswordUser(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error {

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		slog.Error("couldn't find the user id",
			"id", id,
			"error", err.Error())
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		slog.Error("invalid password",
			"oldpassword", oldPassword,
			"newpassword", newPassword,
			"error", err.Error())
		return fmt.Errorf("password doesn't match")
	}

	hashPass, err := hashPassword(newPassword)
	if err != nil {
		slog.Error("invalid password",
			"newpasswrod", newPassword,
			"error", err.Error())
		return err
	}

	user.Password = hashPass

	u.UserRep.DeleteRefreshToken(ctx, id)

	refreshTokenString, err := GenerateRefreshToken(id.String())
	if err != nil {
		slog.Error("Token error when changing password",
			"error", err)
	}

	newRefreshtoken := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.Uuid,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}

	u.UserRep.SaveRefreshToken(ctx, newRefreshtoken)

	return nil

}

func (u *UserService) DeleteAccountUser(ctx context.Context, id uuid.UUID) error {

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		slog.Error("couldn't find the user id",
			"id", id,
			"error", err.Error())
		return err
	}

	return u.UserRep.DeleteUser(ctx, user.Uuid)
}
