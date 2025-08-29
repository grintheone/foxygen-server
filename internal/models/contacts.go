package models

import "github.com/google/uuid"

type Contact struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Position *string   `json:"position" db:"position"`
	Phone    string    `json:"phone" db:"phone"`
	Email    string    `json:"email" db:"email"`
	ClientID string    `json:"client_id,omitempty" db:"client_id"`
}

type ContactUpdate struct {
	Name     *string `json:"name,omitempty" db:"name"`
	Position *string `json:"position,omitempty" db:"position"`
	Phone    *string `json:"phone,omitempty" db:"phone"`
	Email    *string `json:"email,omitempty" db:"email"`
	ClientID *string `json:"client_id,omitempty" db:"client_id"`
}
