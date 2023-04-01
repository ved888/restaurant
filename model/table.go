package model

import (
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type ResTable struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Code      int       `json:"code" db:"code" validate:"required"`
	Capacity  int       `json:"capacity" db:"capacity" validate:"required"`
	BookingId uuid.UUID `json:"bookingId" db:"booking_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt null.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt null.Time `json:"deletedAt" db:"deleted_at"`
}
