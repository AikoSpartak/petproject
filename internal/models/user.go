package models

import (
	"github.com/google/uuid"
)

type User struct {
	Uuid         uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	FullName     string
	Email        string
	PhoneNumber  string
	DateBirthday uint
	Passport     string
	Password     string
	IsBanned     bool
	Role         string
	Country      string
}
