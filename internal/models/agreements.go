package models

import (
	"time"

	"github.com/google/uuid"
)

type Agreement struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Number       *string    `json:"number" db:"number"`
	ActualClient uuid.UUID  `json:"actual_client" db:"actual_client"`
	Distributor  uuid.UUID  `json:"distributor" db:"distributor"`
	Device       *uuid.UUID `json:"device" db:"device"`
	AssignedAt   *time.Time `json:"assigned_at" db:"assigned_at"`
	FinishedAt   *time.Time `json:"finished_at" db:"finished_at"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	OnWarranty   bool       `json:"on_warranty" db:"on_warranty"`
	Type         string     `json:"type" db:"type"`
}

type AgreementCard struct {
	Agreement
	ClientName    string `json:"client_name" db:"client_name"`
	ClientAddress string `json:"client_address" db:"client_address"`
	DeviceName *string `json:"device_name" db:"device_name"`
	DeviceSerial *string `json:"device_serial" db:"device_serial"`
}
