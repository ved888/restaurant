package model

import (
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Address struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Line1     string    `json:"line1" db:"line1" validate:"required" `
	Line2     string    `json:"line2" db:"line2" validate:"required"`
	PinCode   string    `json:"pinCode" db:"pin_code" validate:"required"`
	City      string    `json:"city" db:"city" validate:"required"`
	State     string    `json:"state" db:"state" validate:"required"`
	Country   string    `json:"country" db:"country" validate:"required"`
	UsersId   uuid.UUID `json:"usersId" db:"users_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt null.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt null.Time `json:"deletedAt" db:"deleted_at"`
}
type UserAddress struct {
	Id        uuid.UUID `json:"id" db:"id"`
	UsersId   uuid.UUID `json:"usersId" db:"users_id"`
	AddressId uuid.UUID `json:"addressId" db:"address_id"`
}
