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

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Orders) error

	GetOwnerOrder(ctx context.Context, uuid uuid.UUID) (*models.Orders, error)

	UpdateOrder(ctx context.Context, order *models.Orders) error

	DeleteOrder(ctx context.Context, uuid uuid.UUID) error
}

type OrderRepositoryStruct struct {
	DB    *gorm.DB
	Cache *cache.Cache
}

func (o *OrderRepositoryStruct) CreateOrder(ctx context.Context, order *models.Orders) error {

	if err := o.DB.WithContext(ctx).Create(order).Error; err != nil {
		slog.Error("database error order",
			"error", err.Error())
		return err
	}

	return nil
}

func (o *OrderRepositoryStruct) GetOrder(ctx context.Context, uuid uuid.UUID) (*models.Orders, error) {

	if val, found := o.Cache.Get(uuid.String()); found {
		return val.(*models.Orders), nil
	}

	var order models.Orders

	err := o.DB.WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	o.Cache.Set(uuid.String(), &order, 5*time.Minute)

	return &order, nil
}

func (o *OrderRepositoryStruct) UpdateOrder(ctx context.Context, uuid uuid.UUID, orderUpdate map[string]interface{}) error {

	if err := o.DB.WithContext(ctx).
		Model(&models.Orders{}).
		Where("uuid = ?", uuid.String()).
		Updates(orderUpdate).Error; err != nil {
		slog.Error("failed update orders",
			"error", err.Error())
		return err
	}

	o.Cache.Delete(uuid.String())
	o.Cache.Delete("user:all:*")

	return nil

}

func (o *OrderRepositoryStruct) DeleteOrder(ctx context.Context, uuid uuid.UUID) error {

	if err := o.DB.WithContext(ctx).
		Model(&models.Orders{}).
		Delete("uuid = ?", uuid.String()).Error; err != nil {
		slog.Error("failed delete order",
			"error", err.Error())
	}

	orderKey := fmt.Sprintf("user:%s", uuid.String())
	o.Cache.Delete(orderKey)
	o.Cache.Delete("user:all:*")

	return nil

}
