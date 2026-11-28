package services

import (
	"arch/internal/models"
	"arch/internal/repository"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type OrdersServices interface {
	RegisterOrderServices(ctx context.Context, ownerUuid uuid.UUID, equipmentUuid pq.StringArray, branchUuid uuid.UUID, sellerUuid uuid.UUID, invoicesUuid uuid.UUID, totalAmount decimal.Decimal, isPaid bool) (*models.Orders, error)

	GetOrderServices(ctx context.Context, orderUuid uuid.UUID) (*models.Orders, error)
	GetListOrdersServices(ctx context.Context, filter repository.Orderfilter) ([]*models.Orders, error)

	UpdateOrderServices(ctx context.Context, orderUuid uuid.UUID, updateMap map[string]interface{}) (*models.Orders, error)
	DeleteOrderServices(ctx context.Context, orderuuid uuid.UUID) error
}

type OrdersServicesStruct struct {
	OrdersResotitory *repository.OrderRepositoryStruct
}

func (o *OrdersServicesStruct) RegisterOrderServices(
	ctx context.Context,
	ownerUuid uuid.UUID,
	equipmentUuid pq.StringArray,
	branchUuid uuid.UUID,
	sellerUuid uuid.UUID,
	invoicesUuid uuid.UUID,
	totalAmount decimal.Decimal,
	isPaid bool) (*models.Orders, error) {

	if ownerUuid == uuid.Nil {
		return nil, errors.New("invalid ownerid")
	}

	if branchUuid == uuid.Nil {
		return nil, errors.New("invalid branchuuid")
	}

	if invoicesUuid == uuid.Nil {
		return nil, errors.New("invalid invoicesuuid")
	}

	if totalAmount.Cmp(decimal.Zero) < 0 {
		return nil, errors.New("failed sum")
	}

	if len(equipmentUuid) <= 0 {
		return nil, errors.New("failed sum equipment")

	}

	order := &models.Orders{
		Uuid:          uuid.New(),
		OwnerUuid:     ownerUuid,
		EquipmentUuid: equipmentUuid,
		BranchUuid:    branchUuid,
		SellerUuid:    sellerUuid,
		InvoicesUuid:  invoicesUuid,
		CreatedAt:     time.Now(),
		TotalAmount:   totalAmount,
		IsPaid:        isPaid,
	}

	if err := o.OrdersResotitory.CreateOrder(ctx, order); err != nil {
		slog.Error("failed add order",
			"error", err.Error())
		return nil, errors.New("failed create order")
	}

	return order, nil

}

func (o *OrdersServicesStruct) GetOrderServices(ctx context.Context, orderUuid uuid.UUID) (*models.Orders, error) {

	order, err := o.OrdersResotitory.GetOrder(ctx, orderUuid)
	if err != nil {
		slog.Error("order not found",
			"error", err.Error())
		return nil, err
	}

	if order == nil {
		slog.Error("order not found")
		return nil, errors.New("order not found")
	}

	return order, nil
}

func (o *OrdersServicesStruct) GetListOrdersServices(ctx context.Context, filter repository.Orderfilter) ([]*models.Orders, error) {
	return o.OrdersResotitory.GetListOrder(ctx, filter)
}

func (o *OrdersServicesStruct) UpdateOrderServices(
	ctx context.Context,
	orderUuid uuid.UUID,
	updateMap map[string]interface{}) (*models.Orders, error) {

	order, err := o.OrdersResotitory.GetOrder(ctx, orderUuid)
	if err != nil {
		slog.Error("invalid order id",
			"error", err.Error())
		return nil, err
	}

	if len(updateMap) == 0 {
		return order, nil
	}

	err = o.OrdersResotitory.UpdateOrder(ctx, orderUuid, updateMap)
	if err != nil {
		slog.Error("couldn't update user",
			"error", err.Error())
		return nil, err
	}

	orderUpdate, err := o.OrdersResotitory.GetOrder(ctx, orderUuid)
	if err != nil {
		slog.Error("failed order update",
			"error", err.Error())
		return nil, err
	}

	return orderUpdate, nil
}

func (o *OrdersServicesStruct) DeleteOrderServices(ctx context.Context, orderuuid uuid.UUID) error {

	order, err := o.OrdersResotitory.GetOrder(ctx, orderuuid)
	if err != nil {
		slog.Error("invalid order id",
			"error", err.Error())
		return err
	}

	return o.OrdersResotitory.DeleteOrder(ctx, order.Uuid)
}
