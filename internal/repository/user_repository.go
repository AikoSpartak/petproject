package repository

import (
	// "arch/internal/cache"
	"arch/internal/models"
	"context"
	"fmt"
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
	Cache *cache.Cache // ← работает
}

func (u *UserRepo) CreateUser(ctx context.Context, user *models.User) error {

	if err := u.DB.WithContext(ctx).Create(user).Error; err != nil {
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
		return nil, err
	}

	u.Cache.Set(key, &user, 5*time.Minute)

	return &user, nil
}

func (u *UserRepo) GetAllUsers(ctx context.Context, active bool) ([]*models.User, error) {

	key := fmt.Sprintf("users:all:active:%v", active)

	if val, found := u.Cache.Get(key); found {
		return val.([]*models.User), nil
	}

	var users []*models.User //

	query := u.DB.WithContext(ctx)
	if active {
		query = query.Where("active = ?", true)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	u.Cache.Set(key, users, 5*time.Minute)
	return users, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {

	var user models.User
	if err := u.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) UpdateUser(ctx context.Context, user *models.User) error {

	if err := u.DB.WithContext(ctx).Model(user).Updates(user).Error; err != nil {
		return err
	}

	for k := range u.Cache.Items() {
		if len(k) >= 5 && k[:5] == "user:" || len(k) >= 6 && k[:6] == "users:" {
			u.Cache.Delete(k)
		}
	}

	return nil
}

func (u *UserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {

	if err := u.DB.WithContext(ctx).Delete(&models.User{}, "uuid = ?", id).Error; err != nil {
		return err
	}

	for k := range u.Cache.Items() {
		if len(k) >= 5 && k[:5] == "user:" || len(k) >= 6 && k[:6] == "users:" {
			u.Cache.Delete(k)
		}
	}

	return nil
}
