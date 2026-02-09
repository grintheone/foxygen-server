package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Locations is a slice of Location that implements sql.Scanner and driver.Valuer
type Locations []Location

// Scan implements the sql.Scanner interface
func (l *Locations) Scan(value interface{}) error {
	if value == nil {
		*l = []Location{}
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("unsupported type for Locations")
	}

	if len(data) == 0 || string(data) == "null" {
		*l = []Location{}
		return nil
	}

	return json.Unmarshal(data, l)
}

// Value implements the driver.Valuer interface
func (l Locations) Value() (driver.Value, error) {
	if len(l) == 0 {
		return "[]", nil
	}
	return json.Marshal(l)
}

type Client struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	Title            string         `json:"title" db:"title"`
	Region           *uuid.UUID     `json:"region" db:"region"`
	Address          string         `json:"address" db:"address"`
	LaboratorySystem *uuid.UUID     `json:"laboratory_system" db:"laboratory_system"`
	Location         *Locations     `json:"location" db:"location"`
	Manager          pq.StringArray `json:"manager" db:"manager"`
}

type ClientUpdate struct {
	Title            *string         `json:"title,omitempty" db:"title"`
	Region           *uuid.UUID      `json:"region,omitempty" db:"region"`
	Address          *string         `json:"address,omitempty" db:"address"`
	Location         *[]Location     `json:"location,omitempty" db:"location"`
	LaboratorySystem *uuid.UUID      `json:"laboratory_system,omitempty" db:"laboratory_system"`
	Manager          *pq.StringArray `json:"manager,omitempty" db:"manager"`
}
