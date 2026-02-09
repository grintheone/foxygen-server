package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Manufacturer can be a string, an object, or null
type FlexibleManufacturer struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	// Add other fields from the object here
}

func (m *FlexibleManufacturer) UnmarshalJSON(data []byte) error {
	// Handle 'null'
	if string(data) == "null" {
		return nil
	}

	// Check the first byte to determine the JSON type
	switch data[0] {
	case '"': // It's a string (likely just the UUID)
		var idStr string
		if err := json.Unmarshal(data, &idStr); err != nil {
			return err
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return err
		}
		m.ID = id
	case '{': // It's a full object
		// Use a local type to avoid infinite recursion
		type alias FlexibleManufacturer
		return json.Unmarshal(data, (*alias)(m))
	default:
		return fmt.Errorf("unsupported type for Manufacturer")
	}
	return nil
}

// 1. Define a named type for the slice
type MaintenanceRegulations []JSONB

// 2. Implement driver.Valuer to convert the slice to JSON for the database
func (m MaintenanceRegulations) Value() (driver.Value, error) {
	if m == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(m)
}

// 3. Implement sql.Scanner to convert the database JSON back into the slice
func (m MaintenanceRegulations) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, m)
}

type Classificator struct {
	ID                      uuid.UUID              `json:"id" db:"id"`
	Title                   string                 `json:"title" db:"title"`
	Manufacturer            *uuid.UUID             `json:"manufacturer" db:"manufacturer"`
	ResearchType            *uuid.UUID             `json:"researchType" db:"research_type"`
	RegistrationCertificate JSONB                  `json:"registrationCertificate" db:"registration_certificate"`
	MaintenanceRegulations  MaintenanceRegulations `json:"maintenanceRegulations" db:"maintenance_regulations"`
	Attachments             pq.StringArray         `json:"attachments" db:"attachments"`
	Images                  pq.StringArray         `json:"images" db:"images"`
}

type ClassificatorUpdate struct {
	Title                   *string                `json:"title,omitempty" db:"title"`
	Manufacturer            *uuid.UUID             `json:"manufacturer,omitempty" db:"manufacturer"`
	ResearchType            *uuid.UUID             `json:"research_type,omitempty" db:"research_type"`
	RegistrationCertificate *JSONB                 `json:"registration_certificate,omitempty" db:"registration_certificate"`
	MaintenanceRegulations  MaintenanceRegulations `json:"maintenance_regulations,omitempty" db:"maintenance_regulations"`
	Attachments             *pq.StringArray        `json:"attachments,omitempty" db:"attachments"`
	Images                  *pq.StringArray        `json:"images,omitempty" db:"images"`
}
