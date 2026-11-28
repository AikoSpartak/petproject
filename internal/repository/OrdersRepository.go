package repository

import (
	"arch/internal/models"
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Orders) error

	GetOrder(ctx context.Context, uuid uuid.UUID) (*models.Orders, error)
	GetListOrder(ctx context.Context, filter Orderfilter) ([]*models.Orders, error)

	UpdateOrder(ctx context.Context, uuid uuid.UUID, orderUpdate map[string]interface{}) error

	DeleteOrder(ctx context.Context, uuid uuid.UUID) error
}

type OrderRepositoryStruct struct {
	DB    *gorm.DB
	Cache *cache.Cache
}

type Orderfilter struct {
	OwnerUuid    string
	Branchuuid   string
	SellerUuid   string
	InvoicesUuid string
	Limit        int
	Offset       int
	OrderBy      string
}

const maxLimitOrder = 20

func (o *OrderRepositoryStruct) CreateOrder(ctx context.Context, order *models.Orders) error {

	if err := o.DB.WithContext(ctx).Create(order).Error; err != nil {
		slog.Error("database error order",
			"error", err.Error())
		return err
	}

	return nil
}

func (o *OrderRepositoryStruct) GetOrder(ctx context.Context, uuid uuid.UUID) (*models.Orders, error) {

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

	return &order, nil
}

func (o *OrderRepositoryStruct) GetListOrder(ctx context.Context, filter Orderfilter) ([]*models.Orders, error) {

	if filter.Limit <= 0 || filter.Limit > 20 {
		filter.Limit = maxLimitOrder
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	var orders []*models.Orders

	query := o.DB.WithContext(ctx)

	if filter.OwnerUuid != "" {
		query = query.Where("owner_uuid = ?", filter.OwnerUuid)
		query = query.Preload("Owner")
	}

	if filter.Branchuuid != "" {
		query = query.Where("branch_uuid = ?", filter.Branchuuid)
		query = query.Preload("Branch")
	}

	if filter.SellerUuid != "" {
		query = query.Where("seller_uuid = ?", filter.SellerUuid)
		query = query.Preload("Seller")
	}

	if filter.InvoicesUuid != "" {
		query = query.Where("invoices_uuid = ?", filter.InvoicesUuid)
		query = query.Preload("Invoices")
	}

	query = query.
		Order(filter.OrderBy).
		Limit(filter.Limit).
		Offset(filter.Offset)

	if err := query.Find(&orders).Error; err != nil {
		slog.Error("order not found",
			"error", err.Error(),
			"filter", filter)
		return nil, err
	}

	return orders, nil

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

	return nil

}

func (o *OrderRepositoryStruct) DeleteOrder(ctx context.Context, uuid uuid.UUID) error {

	if err := o.DB.WithContext(ctx).
		Where("uuid = ?", uuid).
		Delete(&models.Orders{}).Error; err != nil {
		slog.Error("failed delete order",
			"error", err.Error())
		return err
	}

	return nil

}
