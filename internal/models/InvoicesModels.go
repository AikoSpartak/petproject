package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Invoices struct {
	Uuid         uuid.UUID `gorm:"column:uuid;type:uuid;primaryKey"`
	CustomerUuid uuid.UUID
	EmployeeUuid uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PaidAt       time.Time
	Amount       decimal.Decimal
	IsRejected   bool
	IsPaid       bool
	Customer     *User `gorm:"foreignKey:CustomerUuid;references:Uuid"`
	Employee     *User `gorm:"foreignKey:EmployeeUuid;references:Uuid"`
}
