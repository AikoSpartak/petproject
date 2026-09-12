package ResponseStruct

import "github.com/google/uuid"

type ResponseRegister struct {
	Uuid        uuid.UUID
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
}

type LoginResponse struct {
	Uuid         string `json:"uuid"`
	Email        string `json:"email"`
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}

type GetPrfileResponse struct {
	Uuid         string `json:"uuid"`
	Email        string `json:"email"`
	FullName     string `json:"fullname"`
	PhoneNumber  string `json:"phonenumber"`
	DateBirthday uint   `json:"datebirthday"`
	Role         string `json:"role"`
	Country      string `json:"country"`
}

type UpdateProfileResponse struct {
	Uuid        uuid.UUID `json:"id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phonenumber"`
	Country     string    `json:"country"`
}

type ChangeEmailResponse struct {
	Uuid  uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}
