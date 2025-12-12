package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type TicketsRepository interface {
	ListAllTickets(ctx context.Context, executorID string) ([]*models.TicketCard, error)
	ListAllDepartmentTickets(ctx context.Context, currentUserID string) ([]*models.TicketCard, error)
	GetTicketByID(ctx context.Context, uuid uuid.UUID) (*models.TicketSinglePage, error)
	DeleteTicketByID(ctx context.Context, uuid uuid.UUID) error
	CreateNewTicket(ctx context.Context, payload models.RawTicket) (*models.RawTicket, error)
	UpdateTicketInfo(ctx context.Context, payload models.TicketUpdates, userID string) error
	GetReasonInfoByID(ctx context.Context, id string) (*models.TicketReason, error)
	GetTicketContactPerson(ctx context.Context, uuid uuid.UUID) (*models.Contact, error)
	// GetClientTicketIDs(ctx context.Context, clientUUID uuid.UUID) ([]*uuid.UUID, error)
	GetTicketsByField(ctx context.Context, field string, fieldUUID uuid.UUID, filters models.TicketFilters) (*models.TicketArchiveResponse, error)
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
    t.assigned_interval,
    t.urgent,
    t.status,
    t.workstarted_at,
    t.workfinished_at,
  	t.description,
    TRIM(CONCAT(ex.first_name, ' ', ex.last_name)) as executor,
    dep.title as department,
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
	LEFT JOIN classificators c ON d.classificator = c.id
	LEFT JOIN clients cl ON t.client = cl.id
	LEFT JOIN ticket_reasons tr on t.reason = tr.id
	LEFT JOIN users ex ON t.executor = ex.user_id
	LEFT JOIN departments dep ON t.department = dep.id
	WHERE executor = $1
	ORDER BY
		CASE
        WHEN (t.assigned_interval->>'end')::TIMESTAMP < NOW() THEN 0  -- Overdue first
        WHEN urgent = TRUE THEN 1    -- Then urgent
        ELSE 2                                   -- Then everything else
    END,
    (t.assigned_interval->>'end')::TIMESTAMP ASC;
	`
	// query := "SELECT * FROM tickets WHERE executor = $1"
	var tickets []*models.TicketCard

	err := r.db.SelectContext(ctx, &tickets, query, executorID)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *ticketsRepository) ListAllDepartmentTickets(ctx context.Context, currentUserID string) ([]*models.TicketCard, error) {
	query := `SELECT department FROM users WHERE user_id = $1`

	var department uuid.UUID

	err := r.db.GetContext(ctx, &department, query, currentUserID)
	if err != nil {
		return nil, err
	}

	fmt.Println(department, "department ID")

	query = `
	SELECT
    t.id,
    t.number,
    t.assigned_interval,
    t.urgent,
    t.status,
    t.workstarted_at,
    t.workfinished_at,
  	t.description,
    TRIM(CONCAT(ex.first_name, ' ', ex.last_name)) as executor,
    dep.title as department,
    d.serial_number AS device_serial_number,
    c.title AS device_classificator_title,
    cl.title as client_name,
    cl.address as client_address,
    tr.title as reason
	FROM tickets t
	LEFT JOIN devices d ON t.device = d.id
	LEFT JOIN classificators c ON d.classificator = c.id
	LEFT JOIN clients cl ON t.client = cl.id
	LEFT JOIN ticket_reasons tr on t.reason = tr.id
	LEFT JOIN users ex ON t.executor = ex.user_id
	LEFT JOIN departments dep ON t.department = dep.id
	WHERE t.department = $1
	ORDER BY
		CASE
        WHEN (t.assigned_interval->>'end')::TIMESTAMP < NOW() THEN 0  -- Overdue first
        WHEN urgent = TRUE THEN 1    -- Then urgent
        ELSE 2                                   -- Then everything else
    END,
    (t.assigned_interval->>'end')::TIMESTAMP ASC;
	`
	var tickets []*models.TicketCard

	err = r.db.SelectContext(ctx, &tickets, query, department)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *ticketsRepository) GetTicketByID(ctx context.Context, uuid uuid.UUID) (*models.TicketSinglePage, error) {
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
			t.assigned_interval,
			t.urgent,
			t.status,
			t.result,
			t.used_materials,
			t.recommendation,
			t.executor,
			TRIM(CONCAT(ex.first_name, ' ', ex.last_name)) as executorName,
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

func (r *ticketsRepository) UpdateTicketInfo(ctx context.Context, updates models.TicketUpdates, userID string) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var exists bool
	err = tx.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM tickets WHERE id = $1)", updates.ID)
	if err != nil {
		return fmt.Errorf("failed to check ticket existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("ticket was not found, updating not possible")
	}

	query := "UPDATE tickets SET "
	var setClauses []string
	args := make(map[string]any)
	args["id"] = updates.ID

	// Check each field and add to query if not nil
	if updates.Status != nil {
		setClauses = append(setClauses, "status = :status")
		args["status"] = updates.Status
	}
	if updates.WorkStartedAt != nil {
		setClauses = append(setClauses, "workstarted_at = :workstarted_at")
		args["workstarted_at"] = updates.WorkStartedAt
	}
	if updates.WorkFinishedAt != nil {
		setClauses = append(setClauses, "workfinished_at = :workfinished_at")
		args["workfinished_at"] = updates.WorkFinishedAt
	}
	if updates.Result != nil {
		setClauses = append(setClauses, "result = :result")
		args["result"] = updates.Result
	}
	if updates.Recommendation != nil {
		setClauses = append(setClauses, "recommendation = :recommendation")
		args["recommendation"] = updates.Recommendation
	}
	if updates.ClosedAt != nil {
		setClauses = append(setClauses, "closed_at = :closed_at")
		args["closed_at"] = updates.ClosedAt
	}
	if updates.AssignedAt != nil {
		setClauses = append(setClauses, "assigned_at = :assigned_at")
		args["assigned_at"] = updates.AssignedAt
	}
	if updates.AssignedBy != nil {
		setClauses = append(setClauses, "assigned_by = :assigned_by")
		args["assigned_by"] = updates.AssignedBy
	}
	if updates.Executor != nil {
		setClauses = append(setClauses, "executor = :executor")
		args["executor"] = updates.Executor
	}
	if updates.Description != nil {
		setClauses = append(setClauses, "description = :description")
		args["description"] = updates.Description
	}
	if updates.AssignedInterval != nil {
		setClauses = append(setClauses, "assigned_interval = :assigned_interval")
		args["assigned_interval"] = updates.AssignedInterval
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = :id"

	result, err := tx.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	// Verify ticket was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for ticket update: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no ticket found with ID %q", updates.ID)
	}

	query = `UPDATE users SET latest_ticket = $1 WHERE user_id = $2`
	result, err = tx.ExecContext(ctx, query, updates.ID, userID)
	if err != nil {
		return fmt.Errorf("failed to update user's latest ticket: %w", err)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for user update: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %q", userID)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
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

func (r *ticketsRepository) GetTicketsByField(ctx context.Context, field string, fieldUUID uuid.UUID, filters models.TicketFilters) (*models.TicketArchiveResponse, error) {
	allowedFields := map[string]bool{
		"client":   true,
		"device":   true,
		"executor": true,
	}

	if !allowedFields[field] {
		return nil, fmt.Errorf("invalid filter field: %s", field)
	}

	query := fmt.Sprintf(`
	SELECT
    t.id,
    t.number,
    t.assigned_interval,
    t.urgent,
    t.status,
    t.result,
    t.created_at,
    t.workstarted_at,
    t.workfinished_at,
    TRIM(CONCAT(ex.first_name, ' ', ex.last_name)) as executor,
    dep.title as department,
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
	LEFT JOIN classificators c ON d.classificator = c.id
	LEFT JOIN clients cl ON t.client = cl.id
	LEFT JOIN ticket_reasons tr on t.reason = tr.id
	LEFT JOIN users ex ON t.executor = ex.user_id
	LEFT JOIN departments dep ON ex.department = dep.id
	WHERE %s = $1
			AND (
				($2 = 'closed' AND t.status = 'closed')
				OR ($2 = 'in-progress' AND t.status IN ('inWork', 'worksDone'))
    			OR ($2 = 'all')
			)
	`, field)

	args := []any{
		fieldUUID,
		filters.Status,
	}
	argPos := 3

	if filters.Reason != nil {
		query += fmt.Sprintf(" AND tr.id = $%d", argPos)
		args = append(args, *filters.Reason)
		argPos++
	}

	if filters.DateStart != nil {
		query += fmt.Sprintf(" AND t.created_at >= $%d", argPos)
		args = append(args, *filters.DateStart)
		argPos++
	}

	if filters.DateEnd != nil {
		query += fmt.Sprintf(" AND t.created_at <= $%d", argPos)
		args = append(args, *filters.DateEnd)
		argPos++
	}

	if filters.DeviceID != nil {
		query += fmt.Sprintf(" AND t.device = $%d", argPos)
		args = append(args, *filters.DeviceID)
		argPos++
	}

	query += " ORDER BY created_at"
	var tickets []*models.TicketCard

	err := r.db.SelectContext(ctx, &tickets, query, args...)
	if err != nil {
		return nil, err
	}

	var filterDates []string
	query = fmt.Sprintf(`
		SELECT DISTINCT created_at FROM
		tickets t
		WHERE %s = $1
			AND ($2 = 'closed' AND t.status = 'closed')
  		OR ($2 = 'in-progress' AND t.status IN ('inWork', 'worksDone'))
    	OR ($2 = 'all')
     	ORDER BY created_at DESC
	`, field)

	err = r.db.SelectContext(ctx, &filterDates, query, fieldUUID, filters.Status)
	if err != nil {
		return nil, err
	}

	var reasons []struct {
		Reason string `db:"reason" json:"reason"`
		Title  string `db:"title" json:"title"`
	}

	query = fmt.Sprintf(`
		SELECT DISTINCT
			t.reason,
		 	tr.title as title
		FROM
		tickets t
		LEFT JOIN ticket_reasons tr ON t.reason = tr.id
		WHERE %s = $1
			AND ($2 = 'closed' AND t.status = 'closed')
  			OR ($2 = 'in-progress' AND t.status IN ('inWork', 'worksDone'))
    		OR ($2 = 'all')
	`, field)

	err = r.db.SelectContext(ctx, &reasons, query, fieldUUID, filters.Status)
	if err != nil {
		return nil, err
	}

	var devices []struct {
		ID          string `db:"id" json:"id"`
		DeviceTitle string `db:"device_title" json:"device_title"`
	}

	query = fmt.Sprintf(`
		SELECT DISTINCT
			t.device as id,
			cl.title as device_title
		FROM
		tickets t
		LEFT JOIN devices d ON t.device = d.id
		LEFT JOIN classificators cl ON d.classificator = cl.id
		WHERE %s = $1
			AND ($2 = 'closed' AND t.status = 'closed')
  			OR ($2 = 'in-progress' AND t.status IN ('inWork', 'worksDone'))
    		OR ($2 = 'all')
	`, field)

	err = r.db.SelectContext(ctx, &devices, query, fieldUUID, filters.Status)
	if err != nil {
		return nil, err
	}

	var response = models.TicketArchiveResponse{
		Filters: make(map[string]any),
	}

	response.Tickets = tickets
	response.Filters["availableDates"] = filterDates
	response.Filters["reasons"] = reasons

	if field == "client" {
		response.Filters["devices"] = devices
	}

	return &response, nil
}
