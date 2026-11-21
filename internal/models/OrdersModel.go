package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Orders struct {
	Uuid          uuid.UUID       `gorm:"column:uuid;type:uuid;primaryKey"`
	OwnerUuid     uuid.UUID       `gorm:"column:owner_uuid;type:uuid"`
	EquipmentUuid pq.StringArray  `gorm:"column:equipment_uuid;type:varchar[];default:'{}'"`
	BranchUuid    uuid.UUID       `gorm:"column:branch_uuid;type:uuid"`
	SellerUuid    uuid.UUID       `gorm:"column:seller_uuid;type:uuid"`
	InvoicesUuid  uuid.UUID       `gorm:"column:invoices_uuid;type:uuid"`
	CreatedAt     time.Time       `gorm:"column:created_at;type:date;default:CURRENT_DATE"`
	TotalAmount   decimal.Decimal `gorm:"column:total_amount;type:numeric(10,2)"`
	IsPaid        bool            `gorm:"column:is_paid;type:boolean;default:false"`
	Branch        *Branches       `gorm:"foreignKey:BranchUuid;references:Uuid"`
	Seller        *User           `gorm:"foreignKey:SellerUuid;references:Uuid"`
	Invoices      *Invoices       `gorm:"foreignKey:InvoicesUuid;references:Uuid"`
	Owner         *User           `gorm:"foreignKey:OwnerUuid;references:Uuid"`
}
