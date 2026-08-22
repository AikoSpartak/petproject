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

	start := time.Now()

	err := emailValidate(login)
	if err != nil {
		slog.Warn("mail validation failed",
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
		slog.Warn("phone validation failed",
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

	slog.Info("the user has been added",
		"id", user.Uuid.String(),
		"email", user.Email,
		"phonenumber", user.PhoneNumber,
		"time", time.Since(start).Milliseconds(),
	)

	return user, nil

}

// добавить JWT токен.
func (u *UserService) LoginUser(ctx context.Context, login, password string) (string, *models.User, error) {

	start := time.Now()

	slog.Info("LoginUser started",
		"email", login,
	)

	user, err := u.UserRep.GetByEmail(ctx, login)
	if err != nil {
		slog.Warn("couldn't get the user",
			"email", login,
			"err", err.Error())
		return "", nil, fmt.Errorf("invalid username")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		slog.Error("invalid password",
			"email", login,
			"error", err.Error())
		return "", nil, fmt.Errorf("invalid password")
	}

	token, err := GenerateToken(
		user.Uuid.String(),
		user.Email,
		user.Role,
	)
	if err != nil {
		slog.Error("failed to generate jwt",
			"email", login,
			"error", err.Error(),
		)
		return "", nil, err
	}

	slog.Info("successful authorization",
		"Uuid", user.Uuid.String(),
		"email", user.Email,
		"time", time.Since(start).Milliseconds(),
	)

	return token, user, nil
}

func (u *UserService) GetProfile(ctx context.Context, id uuid.UUID) (*models.User, error) {

	start := time.Now()

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		slog.Warn("couldn't find the user",
			"id", id,
			"error", err.Error())
		return nil, err
	}

	slog.Info("the user has been found",
		"uuid", user.Uuid.String(),
		"email", user.Email,
		"fullname", user.FullName,
		"phonenumber", user.PhoneNumber,
		"datebirthday", user.DateBirthday,
		"role", user.Role,
		"country", user.Country,
		"time", time.Since(start).Microseconds())

	return user, nil
}

func (u *UserService) UpdateProfile(ctx context.Context, id uuid.UUID, name, phone string) (*models.User, error) {

	start := time.Now()

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		slog.Warn("couldn't find the user",
			"id", id,
			"error", err.Error())
		return nil, ErrUserNotFound
	}

	user.FullName = name
	user.PhoneNumber = phone

	u.UserRep.UpdateUser(ctx, user)

	slog.Info("data changed",
		"uuid", user.Uuid.String(),
		"fullname", user.FullName,
		"email", user.Email,
		"phonenumber", user.PhoneNumber,
		"country", user.Country,
		"time", time.Since(start).Microseconds())

	return user, nil
}

func (u *UserService) ChangeEmail(ctx context.Context, id uuid.UUID, newEmail string) (*models.User, error) {

	start := time.Now()

	err := emailValidate(newEmail)
	if err != nil {
		slog.Warn("mail validation failed",
			"email", newEmail,
			"error", err.Error())
		return nil, ErrorSenEmail
	}

	user, err := u.UserRep.GetByEmail(ctx, newEmail)
	if err != nil {
		slog.Warn("couldn't get the user email",
			"error", err.Error())
		return nil, ErrUserNotFound
	}

	if user != nil && user.Uuid != id {
		return nil, err
	}

	userEmail, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		slog.Warn("couldn't find the user id",
			"id", id,
			"error", err.Error())
		return nil, err
	}

	userEmail.Email = newEmail

	err = u.UserRep.UpdateUser(ctx, userEmail)
	if err != nil {
		slog.Warn("couldn't update data",
			"error", err.Error(),
			"email", newEmail)
		return nil, err
	}

	slog.Info("The password has been updated",
		"uuid", user.Uuid.String(),
		"email", user.Email,
		"time", time.Since(start).Microseconds())

	return userEmail, nil

}

func (u *UserService) ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error {

	start := time.Now()

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		slog.Warn("couldn't find the user id",
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

	slog.Info("password has been updated",
		"time", time.Since(start).Milliseconds())

	return nil

}

func (u *UserService) DeleteAccount(ctx context.Context, id uuid.UUID) error {

	start := time.Now()

	user, err := u.UserRep.GetByUserId(ctx, id)
	if err != nil {
		slog.Warn("couldn't find the user id",
			"id", id,
			"error", err.Error())
		return err
	}

	slog.Info("account deleted",
		"time", time.Since(start).Microseconds())

	return u.UserRep.DeleteUser(ctx, user.Uuid)
}
