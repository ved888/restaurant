package model

import (
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Booking struct {
	Id                uuid.UUID `json:"id" db:"id"`
	BookingDate       time.Time `json:"bookingDate" db:"booking_date" validate:"required"`
	PreAdvanceBooking bool      `json:"preAdvanceBooking" db:"pre_advance_booking" validate:"required" `
	UsersId           uuid.UUID `json:"usersId" db:"users_id"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         null.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt         null.Time `json:"deletedAt" db:"deleted_at"`
}
