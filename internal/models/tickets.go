package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TicketReason struct {
	ID      string `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Past    string `json:"past" db:"past"`
	Present string `json:"present" db:"present"`
	Future  string `json:"future" db:"future"`
}

type Ticket struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	Number         string         `json:"number" db:"number"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	AssignedAt     *time.Time     `json:"assigned_at" db:"assigned_at"`
	WorkStartedAt  *time.Time     `json:"workstarted_at" db:"workstarted_at"`
	WorkFinishedAt *time.Time     `json:"workfinished_at" db:"workfinished_at"`
	ClosedAt       *time.Time     `json:"closed_at" db:"closed_at"`
	Client         uuid.UUID      `json:"client" db:"client"`
	Device         uuid.UUID      `json:"device" db:"device"`
	TicketType     string         `json:"ticket_type" db:"ticket_type"`
	Author         uuid.UUID      `json:"author" db:"author"`
	Department     uuid.UUID      `json:"department" db:"department"`
	AssignedBy     uuid.UUID      `json:"assigned_by" db:"assigned_by"`
	Reason         string         `json:"reason" db:"reason"`
	Description    *string        `json:"description" db:"description"`
	ContactPerson  uuid.UUID      `json:"contact_person" db:"contact_person"`
	Executor       uuid.UUID      `json:"executor" db:"executor"`
	Status         string         `json:"status" db:"status"`
	Result         *string        `json:"result" db:"result"`
	UsedMaterials  pq.StringArray `json:"used_materials" db:"used_materials"`
	Recommendation *string        `json:"recommendation" db:"recommendation"`
	Attachments    pq.StringArray `json:"attachments" db:"attachments"`
}

type TicketUpdates struct {
	Number         *string         `json:"number" db:"number"`
	Client         *uuid.UUID      `json:"client" db:"client"`
	Device         *uuid.UUID      `json:"device" db:"device"`
	TicketType     *string         `json:"ticket_type" db:"ticket_type"`
	Author         *uuid.UUID      `json:"author" db:"author"`
	Department     *uuid.UUID      `json:"department" db:"department"`
	AssignedBy     *uuid.UUID      `json:"assigned_by" db:"assigned_by"`
	Reason         *string         `json:"reason" db:"reason"`
	Description    *string         `json:"description" db:"description"`
	ContactPerson  *uuid.UUID      `json:"contact_person" db:"contact_person"`
	Executor       *uuid.UUID      `json:"executor" db:"executor"`
	Status         *string         `json:"status" db:"status"`
	Result         *string         `json:"result" db:"result"`
	UsedMaterials  *pq.StringArray `json:"used_materials" db:"used_materials"`
	Recommendation *string         `json:"recommendation" db:"recommendation"`
	Attachments    *pq.StringArray `json:"attachments" db:"attachments"`
	AssignedAt     *time.Time      `json:"assigned_at" db:"assigned_at"`
	WorkStartedAt  *time.Time      `json:"workstarted_at" db:"workstarted_at"`
	WorkFinishedAt *time.Time      `json:"workfinished_at" db:"workfinished_at"`
	ClosedAt       *time.Time      `json:"closed_at" db:"closed_at"`
}
