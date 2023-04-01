package model

import (
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Billing struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Type      int       `json:"type" db:"type" validate:"required"`
	Mode      string    `json:"mode" db:"mode" validate:"required"`
	OrdersId  uuid.UUID `json:"ordersId" db:"orders_id"`
	UsersId   uuid.UUID `json:"usersId" db:"users_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt null.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt null.Time `json:"deletedAt" db:"deleted_at"`
}
