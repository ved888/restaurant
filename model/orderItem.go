package model

import (
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type OrderItem struct {
	Id        uuid.UUID `json:"id" db:"db"`
	Price     int       `json:"price" db:"price" validate:"required"`
	Quantity  int       `json:"quantity" db:"quantity" validate:"required"`
	OrderId   uuid.UUID `json:"orderId" db:"order_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt null.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt null.Time `json:"deletedAt" db:"deleted_at"`
}
