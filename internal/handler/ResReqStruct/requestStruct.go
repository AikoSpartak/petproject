package resreqstruct

import "github.com/google/uuid"

type RegisterRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PhoneNumber  string `json:"phonenumber"`
	DateBirthday string `json:"datebirthday"`
	Passport     string `json:"passport"`
	Country      string `json:"country"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetPrfileRequest struct {
	ID uuid.UUID `json:"id"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type ChangeEmailRequest struct {
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}
