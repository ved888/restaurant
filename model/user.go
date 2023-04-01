package model

import (
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"time"
)

type Users struct {
	Id         uuid.UUID `json:"id" db:"id" `
	FirstName  string    `json:"firstName" db:"first_name" validate:"required"`
	MiddleName string    `json:"middleName" db:"middle_name" validate:"required"`
	LastName   string    `json:"lastName" db:"last_name" validate:"required"`
	Phone      string    `json:"phone" db:"phone" validate:"required"`
	EmailId    string    `json:"emailId" db:"email_id" validate:"required,email"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  null.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt  null.Time `json:"deletedAt" db:"deleted_at"`
}

type Interest struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt null.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt null.Time `json:"deletedAt" db:"deleted_at"`
}

type UserInterest struct {
	Id         uuid.UUID `json:"Id" db:"id"`
	InterestId uuid.UUID `json:"InterestId" db:"interest_id"`
	UsersId    uuid.UUID `json:"UsersId" db:"users_id"`
}

type UserInterestRequest struct {
	User         Users        `json:"user"`
	Interest     Interest     `json:"interest"`
	Address      Address      `json:"address"`
	UserInterest UserInterest `json:"-"`
	//InterestUser InterestUser `json:"-"`
}

type InterestUser struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Name   string    `json:"name" db:"name"`
	Type   string    `json:"type" db:"type"`
	UserID uuid.UUID `json:"-" db:"user_id"`
}
