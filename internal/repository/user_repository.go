package repository

import (
	"arch/internal/models"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetByUserId(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetAllUsers(ctx context.Context, active bool, limit, offset int, orderBy string) ([]*models.User, error)
	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, id uuid.UUID) error
}

type UserRepo struct {
	DB    *gorm.DB
	Cache *cache.Cache
}

const ttl = 5 * time.Minute
const maxLimit = 50

func (u *UserRepo) CreateUser(ctx context.Context, user *models.User) error {

	if err := u.DB.WithContext(ctx).Create(user).Error; err != nil {
		slog.Error("Failed creat user",
			"error", err.Error())
		return err
	}

	return nil
}

func (u *UserRepo) GetByUserId(ctx context.Context, id uuid.UUID) (*models.User, error) {

	key := fmt.Sprintf("user:%s", id.String())

	if val, found := u.Cache.Get(key); found {
		return val.(*models.User), nil
	}

	var user models.User
	if err := u.DB.WithContext(ctx).Where("uuid = ?", id).First(&user).Error; err != nil {
		slog.Error("couldn't find the user in db",
			"id", id,
			"error", err.Error())

		return nil, err
	}

	u.Cache.Set(key, &user, ttl)

	return &user, nil
}

func (u *UserRepo) GetAllUsers(ctx context.Context, active bool, limit, offset int, orderBy string) ([]*models.User, error) {

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

	u.Cache.Set(key, users, ttl)

	return users, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {

	key := fmt.Sprintf("user:email:%s", email)
	if val, found := u.Cache.Get(key); found {
		return val.(*models.User), nil
	}

	var user models.User
	if err := u.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		slog.Error("unable to find a user by email in the database",
			"email", email,
			"error", err.Error())
		return nil, err
	}

	u.Cache.Set(key, &user, ttl)

	return &user, nil
}

func (u *UserRepo) UpdateUser(ctx context.Context, user *models.User) error {

	if err := u.DB.WithContext(ctx).Model(user).Updates(user).Error; err != nil {
		slog.Error("failed to update a user in the database",
			"error", err.Error())
		return err
	}

	userKey := fmt.Sprintf("user:%s", user.Uuid.String())
	u.Cache.Delete(userKey)

	u.Cache.Delete("user:all:*")

	return nil
}

func (u *UserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {

	if err := u.DB.WithContext(ctx).Delete(&models.User{}, "uuid = ?", id).Error; err != nil {
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

func (u *UserRepo) SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	if err := u.DB.WithContext(ctx).Create(token).Error; err != nil {
		slog.Error("Error saving token",
			"token", token)
	}
	return nil
}

func (u *UserRepo) DeleteRefreshToken(ctx context.Context, id uuid.UUID) error {
	if err := u.DB.WithContext(ctx).Where("token = ?", id).Delete(&models.RefreshToken{}).Error; err != nil {
		slog.Error("Error delete token",
			"id", id,
			"error", err)
	}
	return nil
}
