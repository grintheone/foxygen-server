package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Contact struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Position string    `json:"position" db:"position"`
	Phone    string    `json:"phone" db:"phone"`
	Email    string    `json:"email" db:"email"`
	ClientID uuid.UUID `json:"client_id,omitempty" db:"client_id"`
}

// Scan implements the sql.Scanner interface
func (c *Contact) Scan(value any) error {
	if value == nil {
		return nil
	}

	// PostgreSQL returns the JSON result as []byte (represented as []uint8)
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed: %T", value)
	}

	return json.Unmarshal(bytes, c)
}

type ContactUpdate struct {
	Name     *string    `json:"name,omitempty" db:"name"`
	Position *string    `json:"position,omitempty" db:"position"`
	Phone    *string    `json:"phone,omitempty" db:"phone"`
	Email    *string    `json:"email,omitempty" db:"email"`
	ClientID *uuid.UUID `json:"client_id,omitempty" db:"client_id"`
}
