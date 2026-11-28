package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Rent struct {
	Uuid          uuid.UUID `gorm:"column:uuid;type:uuid;primaryKey"`
	OwnerUuid     uuid.UUID
	EquipmentUuid pq.StringArray
	BranchUuid    uuid.UUID
	EmployeeUuid  uuid.UUID
	InvoicesUuid  uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	RentalDays    int
	ReturnedAt    *time.Time
	IsReturned    bool
	TotalAmount   decimal.Decimal
	Owner         *User     `gorm:"foreignKey:OwnerUuid;references:Uuid"`
	Branch        *Branches `gorm:"foreignKey:BranchUuid;references:Uuid"`
	Employee      *User     `gorm:"foreignKey:EmployeeUuid;references:Uuid"`
	Invoices      *Invoices `gorm:"foreignKey:InvoicesUuid;references:Uuid"`
}
