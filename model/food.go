package model

import (
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Food struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required" `
	Price       int       `json:"price" db:"price" validate:"required"`
	Type        string    `json:"type" db:"type" validate:"required"`
	OrderItemId uuid.UUID `json:"orderItemId" db:"order_item_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   null.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt   null.Time `json:"deletedAt" db:"deleted_at"`
}
