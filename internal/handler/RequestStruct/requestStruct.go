package RequestStruct

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
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
	Name  string `json:"FullName"`
	Phone string `json:"PhoneNumber"`
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

type OrdersRequest struct {
	EquipmentUuid pq.StringArray  `json:"equipment_uuid"`
	BranchUuid    uuid.UUID       `json:"branch_uuid" `
	SellerUuid    uuid.UUID       `json:"seller_uuid" `
	InvoicesUuid  uuid.UUID       `json:"invoices_uuid" `
	TotalAmount   decimal.Decimal `json:"total_amount"`
	IsPaid        bool            `json:"is_paid"`
}
