package ResponseStruct

import (
	"time"

	"github.com/google/uuid"
)

type ResponseRegister struct {
	Uuid        uuid.UUID
	Email       string `json:"Email"`
	PhoneNumber string `json:"PhoneNumber"`
}

type LoginResponse struct {
	Uuid         string `json:"Uuid"`
	Email        string `json:"Email"`
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}

type GetPrfileResponse struct {
	Uuid         string    `json:"Uuid"`
	Email        string    `json:"Email"`
	FullName     string    `json:"FullName"`
	PhoneNumber  string    `json:"PhoneNumber"`
	DateBirthday time.Time `json:"DateBirthday"`
	Role         string    `json:"Role"`
	Country      string    `json:"Country"`
}

type UpdateProfileResponse struct {
	Uuid        uuid.UUID `json:"Id"`
	FullName    string    `json:"FullName"`
	Email       string    `json:"Email"`
	PhoneNumber string    `json:"PhoneNumber"`
	Country     string    `json:"Country"`
}

type ChangeEmailResponse struct {
	Uuid  uuid.UUID `json:"Id"`
	Email string    `json:"Email"`
}

type ChangePasswordResponse struct {
	Message string `json:"Message"`
}

type DeleteResponse struct {
	Message string `json:"Message"`
}
