package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type TicketsRepository interface {
	ListAllTickets(ctx context.Context, executorID string) ([]*models.TicketCard, error)
	GetTicketByID(ctx context.Context, uuid uuid.UUID) (*models.TicketSinglePage, error)
	DeleteTicketByID(ctx context.Context, uuid uuid.UUID) error
	CreateNewTicket(ctx context.Context, payload models.RawTicket) (*models.RawTicket, error)
	UpdateTicketInfo(ctx context.Context, uuid uuid.UUID, payload models.TicketUpdates) (*models.TicketSinglePage, error)
	GetReasonInfoByID(ctx context.Context, id string) (*models.TicketReason, error)
	GetTicketContactPerson(ctx context.Context, uuid uuid.UUID) (*models.Contact, error)
	// GetClientTicketIDs(ctx context.Context, clientUUID uuid.UUID) ([]*uuid.UUID, error)
	GetTicketsByField(ctx context.Context, field string, fieldUUID uuid.UUID) ([]*models.RawTicket, error)
}

type ticketsRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) *ticketsRepository {
	return &ticketsRepository{db}
}

func (r *ticketsRepository) ListAllTickets(ctx context.Context, executorID string) ([]*models.TicketCard, error) {
	query := `
	SELECT
    t.id,
    t.number,
    t.deadline,
    t.urgent,
    t.status,
    t.workstarted_at,
    t.workfinished_at,
    -- Device fields as individual columns
    d.serial_number AS device_serial_number,
    -- Classificator title
    c.title AS device_classificator_title,
    -- Client fields
    cl.title as client_name,
    cl.address as client_address,
    -- Change reason id to readable name
    tr.title as reason
	FROM tickets t
	LEFT JOIN devices d ON t.device = d.id
	LEFT JOIN classificator c ON d.classificator = c.id
	LEFT JOIN clients cl ON t.client = cl.id
	LEFT JOIN ticket_reasons tr on t.reason = tr.id
	WHERE executor = $1
	ORDER BY
	CASE
        WHEN deadline::TIMESTAMP < NOW() THEN 0  -- Overdue first
        WHEN urgent = TRUE THEN 1    -- Then urgent
        ELSE 2                                   -- Then everything else
    END,
    deadline::TIMESTAMP ASC;  -- Sort by deadline within each group
	`
	// query := "SELECT * FROM tickets WHERE executor = $1"
	var tickets []*models.TicketCard

	err := r.db.SelectContext(ctx, &tickets, query, executorID)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *ticketsRepository) GetTicketByID(ctx context.Context, uuid uuid.UUID) (*models.TicketSinglePage, error) {
	// query := `SELECT * FROM tickets WHERE id = $1`
	query := `
		SELECT
			-- Ticket fields
			t.id,
			t.number,
			t.created_at,
			t.assigned_at,
			t.workstarted_at,
			t.workfinished_at,
			t.closed_at,
			t.deadline,
			t.urgent,
			t.status,
			t.result,
			t.used_materials,
			t.recommendation,
			TRIM(CONCAT(ex.first_name, ' ', ex.last_name)) as executor,
			t.ticket_type,
			t.author,
			dep.title as department,
			TRIM(CONCAT(u.first_name, ' ', u.last_name)) as assigned_by,
			t.description,
			CASE
        		WHEN t.status = 'created' THEN tr.future
          		WHEN t.status = 'assigned' THEN tr.future
            	WHEN t.status = 'inWork' THEN tr.present
             	WHEN t.status = 'worksDone' THEN tr.past
              	WHEN t.status = 'closed' THEN tr.past
              	WHEN t.status = 'cancelled' THEN tr.present
               	ELSE 'Нет данных'
        	END AS reason,
			d.id as device_id,
			d.serial_number as device_serial_number,
			c.title as device_classificator_title,
			cl.id as client_id,
			cl.title as client_name,
			cl.address as client_address,
			con.id as contact_person,
			con.position as contact_position,
			con.name as contact_name,
			con.phone as contact_phone
		FROM tickets t
		LEFT JOIN devices d ON t.device = d.id
		LEFT JOIN classificators c ON d.classificator = c.id
		LEFT JOIN contacts con ON t.contact_person = con.id
		LEFT JOIN clients cl ON t.client = cl.id
		LEFT JOIN ticket_reasons tr ON t.reason = tr.id
		LEFT JOIN users u ON t.assigned_by = u.user_id
		LEFT JOIN users ex ON t.executor = ex.user_id
		LEFT JOIN departments dep ON t.department = dep.id
		WHERE t.id = $1;
	`

	var t models.TicketSinglePage

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

func (r *ticketsRepository) CreateNewTicket(ctx context.Context, payload models.RawTicket) (*models.RawTicket, error) {
	query := `
		INSERT INTO tickets (number, client, device, ticket_type, author, assigned_by, reason, contact_person, executor, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING *
	`

	var ticket models.RawTicket

	err := r.db.GetContext(ctx, &ticket, query, payload.Number, payload.Client, payload.Device, payload.TicketType, payload.Author, payload.AssignedBy, payload.Reason, payload.ContactPerson, payload.Executor, payload.Status)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *ticketsRepository) UpdateTicketInfo(ctx context.Context, uuid uuid.UUID, updates models.TicketUpdates) (*models.TicketSinglePage, error) {
	existing, err := r.GetTicketByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if updates.WorkStartedAt != nil {
		existing.WorkStartedAt = updates.WorkStartedAt
	}
	if updates.WorkFinishedAt != nil {
		existing.WorkFinishedAt = updates.WorkFinishedAt
	}
	if updates.Status != nil {
		existing.Status = *updates.Status
	}
	if updates.Result != nil {
		existing.Result = updates.Result
	}
	if updates.Recommendation != nil {
		existing.Recommendation = updates.Recommendation
	}
	if updates.ClosedAt != nil {
		existing.ClosedAt = updates.ClosedAt
	}

	query := `
		UPDATE tickets
		SET number = :number, workstarted_at = :workstarted_at, workfinished_at = :workfinished_at, status = :status,
		result = :result, recommendation = :recommendation, closed_at = :closed_at
		WHERE id = :id
	`

	_, err = r.db.NamedExecContext(ctx, query, &existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (r *ticketsRepository) GetReasonInfoByID(ctx context.Context, id string) (*models.TicketReason, error) {
	query := `SELECT * FROM ticket_reasons WHERE id = $1`

	var reasonData models.TicketReason

	err := r.db.GetContext(ctx, &reasonData, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &reasonData, nil
}

func (r *ticketsRepository) GetTicketContactPerson(ctx context.Context, uuid uuid.UUID) (*models.Contact, error) {
	query := `SELECT * FROM contacts WHERE id = $1`

	var contact models.Contact

	err := r.db.GetContext(ctx, &contact, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &contact, nil
}

func (r *ticketsRepository) GetTicketsByField(ctx context.Context, field string, fieldUUID uuid.UUID) ([]*models.RawTicket, error) {
	allowedFields := map[string]bool{
		"client":   true,
		"device":   true,
		"executor": true,
	}

	if !allowedFields[field] {
		return nil, fmt.Errorf("invalid filter field: %s", field)
	}

	query := fmt.Sprintf(`SELECT * FROM tickets WHERE %s = $1`, field)

	var tickets []*models.RawTicket

	err := r.db.SelectContext(ctx, &tickets, query, fieldUUID)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}
