package repository

import (
	"arch/internal/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

//go:generate mockgen -source=user_repository.go -destination=mocks/mock.go

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetByUserId(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAllUsers(ctx context.Context, filter models.UsersFiltres) ([]*models.User, error)

	UpdateUser(ctx context.Context, id uuid.UUID, userUpdate map[string]interface{}) error
	UpdateUserEmail(ctx context.Context, user *models.User) error
	UpdateUserPassword(ctx context.Context, user *models.User) error

	DeleteUser(ctx context.Context, id uuid.UUID) error

	SaveRefreshToken(userID uuid.UUID, RefreshToken string)
	SaveAccessToken(userID uuid.UUID, AccessToken string)
	DeleteAccessToken(userID uuid.UUID)
	DeleteRefreshToken(userID uuid.UUID)
}

type UserRepositor struct {
	DB    *gorm.DB
	Cache *cache.Cache
}

const ttlAccess = 15 * time.Minute

const maxLimit = 50

func (u *UserRepositor) CreateUser(ctx context.Context, user *models.User) error {

	if err := u.DB.WithContext(ctx).Create(user).Error; err != nil {
		slog.Error("Database error saving new user",
			"error", err.Error())
		return err
	}

	return nil
}

func (u *UserRepositor) GetByUserId(ctx context.Context, id uuid.UUID) (*models.User, error) {

	key := fmt.Sprintf("user:%s", id.String())

	if val, found := u.Cache.Get(key); found {
		return val.(*models.User), nil
	}

	var user models.User
	if err := u.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		slog.Error("User not found",
			"id", id,
			"error", err.Error())

		return nil, err
	}

	u.Cache.Set(key, &user, ttlAccess)

	return &user, nil
}

func (u *UserRepositor) GetAllUsers(ctx context.Context, filter models.UsersFiltres) ([]*models.User, error) {

	if filter.Limit <= 0 || filter.Limit > 50 {
		filter.Limit = maxLimit
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	key := fmt.Sprintf("users:all:%v", filter)

	if val, found := u.Cache.Get(key); found {
		return val.([]*models.User), nil
	}

	var users []*models.User

	query := u.DB.WithContext(ctx)

	if filter.Active {
		query = query.Where("Active = ?", true)
	}

	if filter.Role != "" {
		query = query.Where("Role = ?", filter.Role)
	}

	if filter.Country != "" {
		query = query.Where("Country = ?", filter.Country)
	}

	query = query.
		Order(filter.OrderBy).
		Limit(filter.Limit).
		Offset(filter.Offset)

	if err := query.Find(&users).Error; err != nil {
		slog.Error("User not found",
			"error", err.Error(),
			"filter", filter)
		return nil, err
	}

	u.Cache.Set(key, users, ttlAccess)

	return users, nil
}

func (u *UserRepositor) GetByEmail(ctx context.Context, email string) (*models.User, error) {

	key := fmt.Sprintf("user:email:%s", email)
	if val, found := u.Cache.Get(key); found {
		return val.(*models.User), nil
	}

	var user models.User
	err := u.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("User not found",
			"email", email,
			"error", err.Error())
		return nil, err
	}

	u.Cache.Set(key, &user, ttlAccess)
	return &user, nil
}

func (u *UserRepositor) UpdateUser(ctx context.Context, id uuid.UUID, userUpdate map[string]interface{}) error {

	if err := u.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id.String()).
		Updates(userUpdate).Error; err != nil {
		slog.Error("failed to update in database")
		return err
	}

	userKey := fmt.Sprintf("user:%s", id.String())
	u.Cache.Delete(userKey)
	u.Cache.Delete("user:all:*")

	return nil
}

func (u *UserRepositor) UpdateUserEmail(ctx context.Context, user *models.User) error {

	if err := u.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.Uuid).
		Update("email", user.Email).Error; err != nil {
		slog.Error("failed to update in database")
		return err
	}

	userKey := fmt.Sprintf("user:%s", user.Uuid.String())
	u.Cache.Delete(fmt.Sprintf("user:email:%s", user.Email))
	u.Cache.Delete(userKey)

	return nil

}

func (u *UserRepositor) UpdateUserPassword(ctx context.Context, user *models.User) error {

	if err := u.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.Uuid).
		Update("password", user.Password).Error; err != nil {
		slog.Error("failed to update in database",
			"error", err.Error())
		return err
	}

	userKey := fmt.Sprintf("user:%s", user.Uuid.String())
	u.Cache.Delete(fmt.Sprintf("user:email:%s", user.Email))
	u.Cache.Delete(userKey)
	u.Cache.Delete("user:all:*")

	return nil
}

func (u *UserRepositor) DeleteUser(ctx context.Context, id uuid.UUID) error {

	if err := u.DB.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error; err != nil {
		slog.Error("failed to find a user by ID in the database",
			"id", id,
			"error", err.Error())
		return err
	}

	userKey := fmt.Sprintf("user:%s", id.String())
	u.Cache.Delete(userKey)

	u.Cache.Delete("user:all:*")

	return nil
}

func (u *UserRepositor) SaveAccessToken(userID uuid.UUID, AccessToken string) {

	key := "access" + userID.String()

	u.Cache.Set(key, AccessToken, ttlAccess)

}

func (u *UserRepositor) SaveRefreshToken(userID uuid.UUID, RefreshToken string) {
	key := "refresh" + userID.String()

	u.Cache.Set(key, RefreshToken, 3*24*time.Hour)

}

func (u *UserRepositor) DeleteAccessToken(userID uuid.UUID) {

	key := "access" + userID.String()

	u.Cache.Delete(key)
}

func (u *UserRepositor) DeleteRefreshToken(userID uuid.UUID) {
	key := "refresh" + userID.String()

	u.Cache.Delete(key)
}
