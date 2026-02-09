package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type JSONB map[string]any

// Implement sql.Scanner for Properties
func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = JSONB{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Properties: expected []byte, got %T", value)
	}

	return json.Unmarshal(bytes, j)
}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type Device struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Classificator   uuid.UUID `json:"classificator" db:"classificator"`
	SerialNumber    string    `json:"serialNumber" db:"serial_number"`
	Properties      JSONB     `json:"properties" db:"properties"`
	ConntectedToLIS bool      `json:"connectedToLis" db:"connected_to_lis"`
	IsUsed          bool      `json:"is_used" db:"is_used"`
}

type DeviceSinglePage struct {
	Device
	Classificator string `json:"classificator" db:"classificator"`
}

type DeviceRemoteOption struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Title string    `json:"title" db:"title"`
}

type DeviceUpdates struct {
	Classificator   *uuid.UUID `json:"classificator" db:"classificator"`
	SerialNumber    *string    `json:"serial_number" db:"serial_number"`
	Properties      *JSONB     `json:"properties" db:"properties"`
	ConntectedToLIS *bool      `json:"connected_to_lis" db:"connected_to_lis"`
	IsUsed          *bool      `json:"is_used" db:"is_used"`
}
