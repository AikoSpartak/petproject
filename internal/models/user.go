package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	ID           uint
	FullName     string
	Email        string `gorm:"unique"`
	PhoneNumber  string
	DateBirthday uint
	Passport     string `gorm:"unique"`
	Password     string
	IsBanned     bool
	Role         string
}
