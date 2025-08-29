package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ClientLocation struct {
	Lat float64 `json:"lat" db:"lat"`
	Lng float64 `json:"lng" db:"lng"`
}

// Scan implements the sql.Scanner interface
func (cl *ClientLocation) Scan(value any) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid type for ClientLocation")
	}

	return json.Unmarshal(bytes, cl)
}

// Value implements the driver.Valuer interface
func (cl ClientLocation) Value() (driver.Value, error) {
	return json.Marshal(cl)
}

type Client struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	Title            string         `json:"title" db:"title"`
	Region           int            `json:"region" db:"region"`
	Address          string         `json:"address" db:"address"`
	Location         ClientLocation `json:"location" db:"location"`
	LaboratorySystem uuid.UUID      `json:"laboratory_system" db:"laboratory_system"`
	Manager          pq.StringArray `json:"manager" db:"manager"`
}

type ClientUpdate struct {
	Title            *string         `json:"title,omitempty" db:"title"`
	Region           *int            `json:"region,omitempty" db:"region"`
	Address          *string         `json:"address,omitempty" db:"address"`
	Location         *ClientLocation `json:"location,omitempty" db:"location"`
	LaboratorySystem *uuid.UUID      `json:"laboratory_system,omitempty" db:"laboratory_system"`
	Manager          *pq.StringArray `json:"manager,omitempty" db:"manager"`
}
