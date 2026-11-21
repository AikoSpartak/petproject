package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Rent struct {
	Uuid          uuid.UUID       `gorm:"column:uuid;type:uuid;primaryKey"`
	OwnerUuid     uuid.UUID       `gorm:"column:owner_uuid;type:uuid"`
	EquipmentUuid pq.StringArray  `gorm:"column:equipment_uuids;type:varchar[];default:'{}'"`
	BranchUuid    uuid.UUID       `gorm:"column:branch_uuid;type:uuid"`
	EmployeeUuid  uuid.UUID       `gorm:"column:employee_uuid;type:uuid"`
	InvoicesUuid  uuid.UUID       `gorm:"column:invoices_uuid;type:uuid"`
	CreatedAt     time.Time       `gorm:"column:created_at;type:date;default:CURRENT_DATE"`
	UpdatedAt     time.Time       `gorm:"column:updated_at;type:date;default:CURRENT_DATE"`
	RentalDays    int             `gorm:"column:rental_days;type:int"`
	ReturnedAt    *time.Time      `gorm:"column:returned_at;type:date"`
	IsReturned    bool            `gorm:"column:is_returned;type:bool;default:false"`
	TotalAmount   decimal.Decimal `gorm:"column:total_amount;type:numeric(10,2)"`
	Owner         *User           `gorm:"foreignKey:OwnerUuid;references:Uuid"`
	Branch        *Branches       `gorm:"foreignKey:BranchUuid;references:Uuid"`
	Employee      *User           `gorm:"foreignKey:EmployeeUuid;references:Uuid"`
	Invoices      *Invoices       `gorm:"foreignKey:InvoicesUuid;references:Uuid"`
}
