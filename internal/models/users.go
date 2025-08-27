package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	FirstName  *string   `json:"first_name" db:"first_name"`
	LastName   *string   `json:"last_name" db:"last_name"`
	Department uuid.UUID `json:"department" db:"department"`
	Email      *string   `json:"email" db:"email"`
	Phone      *int      `json:"phone" db:"phone"`
	Userpic    uuid.UUID `json:"user_pic" db:"user_pic"`
}
