package RequestStruct

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email        string    `json:"Email"`
	Password     string    `json:"Password"`
	PhoneNumber  string    `json:"PhoneNumber"`
	DateBirthday time.Time `json:"DateBirthday"`
	Passport     string    `json:"Passport"`
	Country      string    `json:"Country"`
}

type LoginRequest struct {
	Login    string `json:"Email"`
	Password string `json:"Password"`
}

type GetPrfileRequest struct {
	ID uuid.UUID `json:"Id"`
}

type UpdateProfileRequest struct {
	Name  *string `json:"FullName"`
	Phone *string `json:"PhoneNumber"`
}

type ChangeEmailRequest struct {
	Email string `json:"Email"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"OldPassword"`
	NewPassword string `json:"NewPassword"`
}

type DeleteUserRequest struct {
	Email string `json:"Email"`
}
