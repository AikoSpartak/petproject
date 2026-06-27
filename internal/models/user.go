package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	ID                uint
	Full_name         string
	Email             string
	Phone_number      uint
	Date_birthday     uint
	Serial_passport   string
	Number_passport   string
	Sistory_operation *string
	History_refund    *string
	History_moving    *string
}
