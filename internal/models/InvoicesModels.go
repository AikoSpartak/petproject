package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Invoices struct {
	Uuid         uuid.UUID       `gorm:"column:uuid;type:uuid;primaryKey"`
	CustomerUuid uuid.UUID       `gorm:"column:customer_uuid;type:uuid"`
	EmployeeUuid uuid.UUID       `gorm:"column:employee_uuid;type:uuid"`
	CreatedAt    *time.Time      `gorm:"column:created_at;type:date;default:CURRENT_DATE"`
	UpdatedAt    *time.Time      `gorm:"column:updated_at;type:date;default:CURRENT_DATE"`
	PaidAt       *time.Time      `gorm:"column:paid_at;type:date"`
	Amount       decimal.Decimal `gorm:"column:amount;type:numeric(10,2)"`
	IsRejected   bool            `gorm:"column:is_rejected;type:bool;default:false"`
	IsPaid       bool            `gorm:"column:is_paid;type:bool;default:false"`
	Customer     *User           `gorm:"foreignKey:CustomerUuid;references:Uuid"`
	Employee     *User           `gorm:"foreignKey:EmployeeUuid;references:Uuid"`
}
