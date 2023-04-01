package model

import (
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Orders struct {
	Id           uuid.UUID `json:"id" db:"id"`
	ItemDiscount int       `json:"itemDiscount" db:"item_discount" validate:"required"`
	Tax          int       `json:"tax" db:"tax" validate:"required"`
	Shipping     string    `json:"shipping" db:"shipping" validate:"required"`
	Total        int       `json:"total" db:"total" validate:"required"`
	UsersId      uuid.UUID `json:"usersId" db:"users_id"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    null.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt    null.Time `json:"deletedAt" db:"deleted_at"`
}
