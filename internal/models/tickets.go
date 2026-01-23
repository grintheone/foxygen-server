package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TicketReason struct {
	ID      string `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Past    string `json:"past,omitempty" db:"past"`
	Present string `json:"present,omitempty" db:"present"`
	Future  string `json:"future,omitempty" db:"future"`
}

type TicketSinglePage struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	Number           string         `json:"number" db:"number"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	AssignedAt       *time.Time     `json:"assigned_at" db:"assigned_at"`
	WorkStartedAt    *time.Time     `json:"workstarted_at" db:"workstarted_at"`
	WorkFinishedAt   *time.Time     `json:"workfinished_at" db:"workfinished_at"`
	ClosedAt         *time.Time     `json:"closed_at" db:"closed_at"`
	AssignedInterval JSONB          `json:"assigned_interval" db:"assigned_interval"`
	Urgent           bool           `json:"urgent" db:"urgent"`
	Executor         *string        `json:"executor" db:"executor"`
	ExecutorName     string         `json:"executorName"`
	Status           string         `json:"status" db:"status"`
	Result           *string        `json:"result" db:"result"`
	UsedMaterials    pq.StringArray `json:"used_materials" db:"used_materials"`
	TicketType       string         `json:"ticket_type" db:"ticket_type"`
	Author           uuid.UUID      `json:"author" db:"author"`
	Department       string         `json:"department" db:"department"`
	AssignedBy       string         `json:"assigned_by" db:"assigned_by"`
	Reason           string         `json:"reason" db:"reason"`
	Description      *string        `json:"description" db:"description"`
	// Client fields
	ClientID      string  `json:"client_id" db:"client_id"`
	ClientName    *string `json:"client_name" db:"client_name"`
	ClientAddress *string `json:"client_address" db:"client_address"`
	// Device fields
	DeviceID                 string  `json:"device_id" db:"device_id"`
	DeviceSerialNumber       *string `json:"device_serial_number" db:"device_serial_number"`
	DeviceClassificatorTitle *string `json:"device_classificator_title" db:"device_classificator_title"`
	// Contact
	ContactPerson   uuid.UUID `json:"contact_person" db:"contact_person"`
	ContactName     string    `json:"contact_name" db:"contact_name"`
	ContactPosition string    `json:"contact_position" db:"contact_position"`
	ContactPhone    string    `json:"contact_phone" db:"contact_phone"`
}

type RawTicket struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	Number           string         `json:"number" db:"number"`
	CreatedAt        *time.Time     `json:"created_at,omitempty" db:"created_at"`
	AssignedAt       *time.Time     `json:"assigned_at" db:"assigned_at"`
	WorkStartedAt    *time.Time     `json:"workstarted_at" db:"workstarted_at"`
	WorkFinishedAt   *time.Time     `json:"workfinished_at" db:"workfinished_at"`
	ClosedAt         *time.Time     `json:"closed_at" db:"closed_at"`
	AssignedInterval JSONB          `json:"assigned_interval" db:"assigned_interval"`
	Urgent           bool           `json:"urgent" db:"urgent"`
	Executor         uuid.UUID      `json:"executor" db:"executor"`
	Status           string         `json:"status" db:"status"`
	Result           *string        `json:"result" db:"result"`
	UsedMaterials    pq.StringArray `json:"used_materials" db:"used_materials"`
	Recommendation   *string        `json:"recommendation" db:"recommendation"`
	TicketType       string         `json:"ticket_type" db:"ticket_type"`
	Author           uuid.UUID      `json:"author" db:"author"`
	Department       uuid.UUID      `json:"department" db:"department"`
	AssignedBy       uuid.UUID      `json:"assigned_by" db:"assigned_by"`
	Reason           string         `json:"reason" db:"reason"`
	Description      string         `json:"description" db:"description"`
	Client           uuid.UUID      `json:"client" db:"client"`
	Device           uuid.UUID      `json:"device" db:"device"`
	ContactPerson    *uuid.UUID     `json:"contact_person,omitempty" db:"contact_person,omitempty"`
	ReferenceTicket  uuid.UUID      `json:"reference_ticket" db:"reference_ticket"`
}

type TicketCard struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	Number           string     `json:"number" db:"number"`
	AssignedInterval JSONB      `json:"assigned_interval" db:"assigned_interval"`
	Urgent           bool       `json:"urgent" db:"urgent"`
	Reason           string     `json:"reason" db:"reason"`
	Status           string     `json:"status" db:"status"`
	Result           *string    `json:"result" db:"result"`
	Description      *string    `json:"description" db:"description"`
	WorkStartedAt    *time.Time `json:"workstarted_at" db:"workstarted_at"`
	WorkFinishedAt   *time.Time `json:"workfinished_at" db:"workfinished_at"`
	Executor         string     `json:"executor" db:"executor"`
	Department       string     `json:"department" db:"department"`
	// Device fields
	DeviceSerialNumber       *string `json:"device_serial_number" db:"device_serial_number"`
	DeviceClassificatorTitle *string `json:"device_classificator_title" db:"device_classificator_title"`
	// Client fields
	ClientName    *string   `json:"client_name" db:"client_name"`
	ClientAddress *string   `json:"client_address" db:"client_address"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type TicketUpdates struct {
	ID               uuid.UUID  `db:"id"`
	Status           *string    `json:"status,omitempty" db:"status"`
	WorkStartedAt    *time.Time `json:"workstarted_at,omitempty" db:"workstarted_at"`
	WorkFinishedAt   *time.Time `json:"workfinished_at,omitempty" db:"workfinished_at"`
	Result           *string    `json:"result,omitempty" db:"result"`
	Recommendation   *string    `json:"recommendation,omitempty" db:"recommendation"`
	Department       *uuid.UUID `json:"department" db:"department"`
	ClosedAt         *time.Time `json:"closed_at,omitempty" db:"closed_at"`
	AssignedAt       *time.Time `json:"assigned_at,omitempty" db:"assigned_at"`
	AssignedBy       *string    `json:"assigned_by,omitempty" db:"assigned_by"`
	Executor         *uuid.UUID `json:"executor,omitempty" db:"executor"`
	Description      *string    `json:"description,omitempty" db:"description"`
	AssignedInterval *JSONB     `json:"assigned_interval,omitempty" db:"assigned_interval"`
}

type CloseTicket struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	Result         string     `json:"result" db:"result"`
	ClosedAt       string     `json:"closed_at" db:"closed_at"`
	Recommendation *string    `json:"recommendation" db:"recommendation"`
	Department     *uuid.UUID `json:"department" db:"department"`
	DoubleSigned   bool       `json:"double_signed" db:"double_signed"`
}

type TicketFilters struct {
	Department string     `json:"department"`
	Status     string     `json:"status"`
	GroupBy    *string    `json:"groupBy,omitempty"`
	Reason     *string    `json:"reason,omitempty"`
	DateStart  *time.Time `json:"dateStart,omitempty"`
	DateEnd    *time.Time `json:"dateEnd,omitempty"`
	DeviceID   *uuid.UUID `json:"deviceID,omitempty"`
}

type TicketArchiveResponse struct {
	Tickets []*TicketCard  `json:"tickets"`
	Filters map[string]any `json:"filters"`
}

type TicketArchiveResponseGrouped struct {
	TicketArchiveResponse
	Tickets map[string][]*TicketCard `json:"tickets"`
}
