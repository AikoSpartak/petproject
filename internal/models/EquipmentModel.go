package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Status = string

const (
	StatusInStock       Status = "На складе"
	StatusInRental      Status = "В прокате"
	StatusBuyBack       Status = "Продан с обратным выкупом"
	StatusReturned      Status = "Возвращен"
	StatusRepair        Status = "В ремонте"
	StatusWrittenOff    Status = "Списан"
	StatusMovesBranches Status = "Перемещается между фелиалами"
)

type Equipment struct {
	Uuid          uuid.UUID
	Status        Status
	OwnerUuid     uuid.UUID
	BranchUuid    uuid.UUID
	Description   *string
	Category      string
	Brand         string
	Model         string
	Specification json.RawMessage
	SerialNumber  string
	Price         float64
	RendRate      json.RawMessage
	Owner         *User     `gorm:"foreignKey:OwnerUuid;references:Uuid"`
	Branch        *Branches `gorm:"foreignKey:BranchUuid;references:Uuid"`
}
