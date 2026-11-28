package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Orders struct {
	Uuid          uuid.UUID `gorm:"column:uuid;type:uuid;primaryKey"`
	OwnerUuid     uuid.UUID
	EquipmentUuid pq.StringArray `gorm:"column:equipment_uuid;type:varchar[]"`
	BranchUuid    uuid.UUID
	SellerUuid    uuid.UUID
	InvoicesUuid  uuid.UUID
	CreatedAt     time.Time
	TotalAmount   decimal.Decimal
	IsPaid        bool
	Branch        *Branches `gorm:"foreignKey:BranchUuid;references:Uuid"`
	Seller        *User     `gorm:"foreignKey:SellerUuid;references:Uuid"`
	Invoices      *Invoices `gorm:"foreignKey:InvoicesUuid;references:Uuid"`
	Owner         *User     `gorm:"foreignKey:OwnerUuid;references:Uuid"`
}
