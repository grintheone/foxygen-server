package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type TicketsRepository interface {
	ListAllTickets(ctx context.Context) (*[]models.Ticket, error)
	GetTicketByID(ctx context.Context, uuid uuid.UUID) (*models.Ticket, error)
	DeleteTicketByID(ctx context.Context, uuid uuid.UUID) error
	CreateNewTicket(ctx context.Context, payload models.Ticket) (*models.Ticket, error)
	UpdateTicketInfo(ctx context.Context, uuid uuid.UUID, payload models.TicketUpdates) (*models.Ticket, error)
}

type ticketsRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) *ticketsRepository {
	return &ticketsRepository{db}
}

func (r *ticketsRepository) ListAllTickets(ctx context.Context) (*[]models.Ticket, error) {
	var tickets []models.Ticket

	err := r.db.SelectContext(ctx, &tickets, `SELECT * FROM tickets;`)
	if err != nil {
		return nil, err
	}

	return &tickets, nil
}

func (r *ticketsRepository) GetTicketByID(ctx context.Context, uuid uuid.UUID) (*models.Ticket, error) {
	query := `SELECT * FROM tickets WHERE id = $1`
	var t models.Ticket

	err := r.db.GetContext(ctx, &t, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &t, nil
}

func (r *ticketsRepository) DeleteTicketByID(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM tickets WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (r *ticketsRepository) CreateNewTicket(ctx context.Context, payload models.Ticket) (*models.Ticket, error) {
	query := `
		INSERT INTO tickets (number, client, device, ticket_type, author, assigned_by, reason, contact_person, executor, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING *
	`

	var ticket models.Ticket

	err := r.db.GetContext(ctx, &ticket, query, payload.Number, payload.Client, payload.Device, payload.TicketType, payload.Author, payload.AssignedBy, payload.Reason, payload.ContactPerson, payload.Executor, payload.Status)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *ticketsRepository) UpdateTicketInfo(ctx context.Context, uuid uuid.UUID, payload models.TicketUpdates) (*models.Ticket, error) {
	existing, err := r.GetTicketByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if payload.Number != nil {
		existing.Number = *payload.Number
	}
	if payload.Client != nil {
		existing.Client = *payload.Client
	}
	if payload.Device != nil {
		existing.Device = *payload.Device
	}
	if payload.TicketType != nil {
		existing.TicketType = *payload.TicketType
	}
	if payload.Author != nil {
		existing.Author = *payload.Author
	}
	if payload.PlannedInterval != nil {
		existing.PlannedInterval = *payload.PlannedInterval
	}
	if payload.AssignedInterval != nil {
		existing.AssignedInterval = *payload.AssignedInterval
	}
	if payload.ActualInterval != nil {
		existing.ActualInterval = *payload.ActualInterval
	}
	if payload.Department != nil {
		existing.Department = *payload.Department
	}
	if payload.AssignedBy != nil {
		existing.AssignedBy = *payload.AssignedBy
	}
	if payload.AssignedAt != nil {
		existing.AssignedAt = *payload.AssignedAt
	}
	if payload.Reason != nil {
		existing.Reason = *payload.Reason
	}
	if payload.Description != nil {
		existing.Description = payload.Description
	}
	if payload.ContactPerson != nil {
		existing.ContactPerson = *payload.ContactPerson
	}
	if payload.Executor != nil {
		existing.Executor = *payload.Executor
	}
	if payload.Status != nil {
		existing.Status = *payload.Status
	}
	if payload.Result != nil {
		existing.Result = payload.Result
	}
	if payload.UsedMaterials != nil {
		existing.UsedMaterials = *payload.UsedMaterials
	}
	if payload.Recommendation != nil {
		existing.Recommendation = payload.Recommendation
	}
	if payload.Attachments != nil {
		existing.Attachments = *payload.Attachments
	}
	if payload.ClosedAt != nil {
		existing.ClosedAt = payload.ClosedAt
	}

	query := `
		UPDATE tickets
		SET number = :number, client = :client, device = :device, ticket_type = :ticket_type, author = :author, planned_interval = :planned_interval, assigned_interval = :assigned_interval, actual_interval = :actual_interval, department = :department, assigned_by = :assigned_by, assigned_at = :assigned_at, reason = :reason, description = :description, contact_person = :contact_person, executor = :executor, status = :status, result = :result, used_materials = :used_materials, recommendation = :recommendation, attachments = :attachments, closed_at = :closed_at
		WHERE id = :id
	`

	_, err = r.db.NamedExecContext(ctx, query, &existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}
