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
	GetAllUsers(ctx context.Context, active bool) ([]*models.User, error)
}

type UserRepo struct {
	DB    *gorm.DB
	Cache *cache.Cache
}

func (u *UserRepo) CreateUser(ctx context.Context, user *models.User) error {

	start := time.Now()

	slog.Info("Creating a user")

	if err := u.DB.WithContext(ctx).Create(user).Error; err != nil {
		slog.Error("Failed creat user",
			"error", err.Error())
		return err
	}

	slog.Info("Creating completed",
		"time", time.Since(start).Milliseconds())

	return nil
}

func (u *UserRepo) GetByUserId(ctx context.Context, id uuid.UUID) (*models.User, error) {

	start := time.Now()

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

	u.Cache.Set(key, &user, 5*time.Minute)

	slog.Info("the user was found in db",
		"id", user.Uuid.String(),
		"time", time.Since(start).Milliseconds())

	return &user, nil
}

func (u *UserRepo) GetAllUsers(ctx context.Context, active bool) ([]*models.User, error) {

	start := time.Now()

	key := fmt.Sprintf("users:all:active:%v", active)

	if val, found := u.Cache.Get(key); found {
		return val.([]*models.User), nil
	}

	var users []*models.User

	query := u.DB.WithContext(ctx)
	if active {
		query = query.Where("active = ?", true)
	}

	if err := query.Find(&users).Error; err != nil {
		slog.Error("couldn't find the user in db",
			"error", err.Error())
		return nil, err
	}

	u.Cache.Set(key, users, 5*time.Minute)

	slog.Info("the user was found in db",
		"time", time.Since(start).Milliseconds())

	return users, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {

	start := time.Now()

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

	slog.Info("The user was found by mail in db",
		"time", time.Since(start).Microseconds())

	u.Cache.Set(key, &user, 5*time.Minute)

	return &user, nil
}

func (u *UserRepo) UpdateUser(ctx context.Context, user *models.User) error {

	start := time.Now()

	if err := u.DB.WithContext(ctx).Model(user).Updates(user).Error; err != nil {
		slog.Error("failed to update a user in the database",
			"error", err.Error())
		return err
	}

	userKey := fmt.Sprintf("user:%s", user.Uuid.String())
	u.Cache.Delete(userKey)

	u.Cache.Delete("user:all:active:true")
	u.Cache.Delete("user:all:active:false")

	slog.Info("user updated in db",
		"time", time.Since(start).Milliseconds())

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

	u.Cache.Delete("user:all:active:true")
	u.Cache.Delete("user:all:active:false")

	slog.Info("the user is deleted and db")

	return nil
}
