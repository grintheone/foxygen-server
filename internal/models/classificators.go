package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Classificator struct {
	ID                      uuid.UUID      `json:"id" db:"id"`
	Title                   string         `json:"title" db:"title"`
	Manufacturer            *uuid.UUID     `json:"manufacturer" db:"manufacturer"`
	ResearchType            *uuid.UUID     `json:"research_type" db:"research_type"`
	RegistrationCertificate JSONB          `json:"registration_certificate" db:"registration_certificate"`
	MaintenanceRegulations  JSONB          `json:"maintenance_regulations" db:"maintenance_regulations"`
	Attachments             pq.StringArray `json:"attachments" db:"attachments"`
	Images                  pq.StringArray `json:"images" db:"images"`
}

type ClassificatorUpdate struct {
	Title                   *string         `json:"title,omitempty" db:"title"`
	Manufacturer            *uuid.UUID      `json:"manufacturer,omitempty" db:"manufacturer"`
	ResearchType            *uuid.UUID      `json:"research_type,omitempty" db:"research_type"`
	RegistrationCertificate *JSONB          `json:"registration_certificate,omitempty" db:"registration_certificate"`
	MaintenanceRegulations  *JSONB          `json:"maintenance_regulations,omitempty" db:"maintenance_regulations"`
	Attachments             *pq.StringArray `json:"attachments,omitempty" db:"attachments"`
	Images                  *pq.StringArray `json:"images,omitempty" db:"images"`
}
