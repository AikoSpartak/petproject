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
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetAllUsers(ctx context.Context, active bool, limit, offset int, orderBy string) ([]*models.User, error)

	SaveRefreshToken(userID uuid.UUID, RefreshToken string)
	SaveAccessToken(userID uuid.UUID, AccessToken string)
	DeleteAccessToken(userID uuid.UUID)
	DeleteRefreshToken(userID uuid.UUID)
}

type UserRepositor struct {
	DB    *gorm.DB
	Cache *cache.Cache
}

const ttlAccess = 5 * time.Minute

const maxLimit = 50

func (u *UserRepositor) CreateUser(ctx context.Context, user *models.User) error {

	if err := u.DB.WithContext(ctx).Create(user).Error; err != nil {
		slog.Error("Failed creat user",
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
		slog.Error("couldn't find the user in db",
			"id", id,
			"error", err.Error())

		return nil, err
	}

	u.Cache.Set(key, &user, ttlAccess)

	return &user, nil
}

func (u *UserRepositor) GetAllUsers(ctx context.Context, active bool, limit, offset int, orderBy string) ([]*models.User, error) {

	if limit <= 0 || limit > 50 {
		limit = maxLimit
	}

	if offset < 0 {
		offset = 0
	}

	key := fmt.Sprintf("users:all:active:%v:limit:%d:offset:%d:orderBy:%s", active, limit, offset, orderBy)

	if val, found := u.Cache.Get(key); found {
		return val.([]*models.User), nil
	}

	var users []*models.User

	query := u.DB.WithContext(ctx)
	if active {
		query = query.Where("active = ?", true)
	}
	query = query.
		Order(orderBy).
		Limit(limit).
		Offset(offset)

	if err := query.Find(&users).Error; err != nil {
		slog.Error("couldn't find the user in db",
			"error", err.Error(),
			"active", active,
			"limit", limit,
			"offset", offset,
			"orderby", orderBy)
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
		slog.Error("unable user in database",
			"email", email,
			"error", err.Error())
		return nil, err
	}

	u.Cache.Set(key, &user, ttlAccess)
	return &user, nil
}

func (u *UserRepositor) UpdateUser(ctx context.Context, user *models.User) error {

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
