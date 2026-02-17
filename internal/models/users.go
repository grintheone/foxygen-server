package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type JSONBMap map[string]any

func (j *JSONBMap) Scan(value any) error {
	if value == nil {
		*j = nil // Set to nil pointer, not empty map
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONBMap", value)
	}

	// If it's a SQL NULL literal, handle it
	if string(bytes) == "null" {
		*j = nil
		return nil
	}

	var result map[string]any
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = result
	return nil
}

func (j *JSONBMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// type User struct {
// 	UserID          uuid.UUID `json:"user_id" db:"user_id"`
// 	FirstName       *string   `json:"first_name" db:"first_name"`
// 	LastName        *string   `json:"last_name" db:"last_name"`
// 	Department      uuid.UUID `json:"department" db:"department"`
// 	DepartmentTitle string    `json:"department_title" db:"department_title"`
// 	Email           *string   `json:"email" db:"email"`
// 	Phone           *int      `json:"phone" db:"phone"`
// 	Userpic         uuid.UUID `json:"user_pic" db:"user_pic"`
// Properties      *JSONBMap `json:"properties" db:"properties"` // INSTEAD OF latest_ticket field from DB
// ActiveTickets   *int      `json:"active_tickets" db:"active_tickets"`
// }

type User struct {
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	FirstName     string     `json:"firstName" db:"first_name"`
	LastName      string     `json:"lastName" db:"last_name"`
	Department    *uuid.UUID `json:"department" db:"department"`
	Email         string     `json:"email" db:"email"`
	Phone         string     `json:"phone" db:"phone"`
	Logo          string     `json:"logo" db:"logo"`
	Properties    *JSONBMap  `json:"properties" db:"properties"` // INSTEAD OF latest_ticket field from DB
	ActiveTickets *int       `json:"active_tickets" db:"active_tickets"`
	LatestTicket  *uuid.UUID `json:"latest_ticket" db:"latest_ticket"`
}

type UserProfile struct {
	ID                uuid.UUID     `json:"user_id" db:"user_id"`
	Fullname          string        `json:"fullname" db:"fullname"`
	Department        string        `json:"department" db:"department"`
	Email             *string       `json:"email" db:"email"`
	Phone             *string       `json:"phone" db:"phone"`
	Logo              *string       `json:"logo" db:"logo"`
	Properties        *JSONBMap     `json:"properties" db:"properties"`
	ActiveTickets     []*TicketCard `json:"active_tickets" db:"active_tickets"`
	ClosedTickets     *int          `json:"closed_tickets" db:"closed_tickets"`
	Overdue           *int          `json:"overdue" db:"overdue"`
	TicketsInProgress *int          `json:"tickets_in_progress" db:"tickets_in_progress"`
}
