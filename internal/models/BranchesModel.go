package models

import (
	"github.com/google/uuid"
)

type Branches struct {
	Uuid        uuid.UUID `gorm:"column:uuid;type:uuid;primaryKey"`
	Address     string
	Description string
}
