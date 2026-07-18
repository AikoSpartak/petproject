package models

import (
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Uuid         uuid.UUID
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
