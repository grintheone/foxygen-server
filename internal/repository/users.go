package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type UsersRepository interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	ListDepartmentUsers(ctx context.Context, userID string) ([]*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	ListUsers(ctx context.Context) (*[]models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	UpdateUser(ctx context.Context, user models.User) error
}

type usersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) UsersRepository {
	return &usersRepository{db}
}

func (r *usersRepository) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	query := `	
				SELECT 
        u.user_id,
  			CONCAT(u.first_name, ' ',  u.last_name) as fullname,
        dep.title as department,
        u.email,
        u.phone,
        u.user_pic,
        CASE 
            WHEN u.latest_ticket IS NULL THEN '{}'::jsonb
            ELSE jsonb_build_object(
                'status', t.status,
                'workstarted_at', t.workstarted_at,
                'workfinished_at', t.workfinished_at,
                'client_name', c.title
            )
        END AS properties,
          (
            SELECT COUNT(*) 
            FROM tickets 
            WHERE executor = u.user_id 
            AND status = 'closed' 
        ) as closed_tickets,
        (
            SELECT COUNT(*) as overdue_tickets_count
            FROM tickets
            WHERE 
            executor = u.user_id
            AND (assigned_interval->>'end')::timestamp < NOW()
            AND status NOT IN ('closed', 'cancelled')
        ) as overdue,
        (
            SELECT COUNT(*) as tickets_in_progress
            FROM tickets
            WHERE executor = u.user_id 
            AND status IN ('assigned', 'inWork', 'worksDone')
        )
        FROM 
            users u
        LEFT JOIN tickets t on u.latest_ticket = t.id
        LEFT JOIN clients c on t.client = c.id
  			LEFT JOIN departments dep ON u.department = dep.id
        WHERE 
            u.user_id = $1
	`

	var profile models.UserProfile
	err := r.db.GetContext(ctx, &profile, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	query = `
	SELECT
    t.id,
    t.number,
    t.assigned_interval,
    t.urgent,
    t.status,
    t.workstarted_at,
    t.workfinished_at,
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
	WHERE executor = $1
	AND status NOT IN ('created', 'cancelled', 'closed')
	ORDER BY
		CASE
        WHEN (t.assigned_interval->>'end')::TIMESTAMP < NOW() THEN 0  -- Overdue first
        WHEN urgent = TRUE THEN 1    -- Then urgent
        ELSE 2                                   -- Then everything else
    END,
    (t.assigned_interval->>'end')::TIMESTAMP ASC;
	`

	var activeTickets []*models.TicketCard

	err = r.db.SelectContext(ctx, &activeTickets, query, userID)
	if err != nil {
		return nil, err
	}

	profile.ActiveTickets = activeTickets

	return &profile, nil
}

func (r *usersRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `
						SELECT 
            u.user_id,
            u.first_name,
            u.last_name,
            u.department,
  					dep.title as department_title,
            u.email,
            u.phone,
            u.user_pic,
            CASE 
                WHEN u.latest_ticket IS NULL THEN '{}'::jsonb
                ELSE jsonb_build_object(
                    'status', t.status,
                    'workstarted_at', t.workstarted_at,
                    'workfinished_at', t.workfinished_at,
                    'client_name', c.title
                )
            END AS properties,
            (
                SELECT COUNT(*) 
                FROM tickets 
                WHERE executor = u.user_id 
                AND status not in ('closed', 'cancelled')
            ) as active_tickets
        FROM 
            users u
        LEFT JOIN tickets t on u.latest_ticket = t.id
        LEFT JOIN clients c on t.client = c.id
  			LEFT JOIN departments dep ON u.department = dep.id
        where 
            u.user_id = $1
`

	var user models.User

	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *usersRepository) ListUsers(ctx context.Context) (*[]models.User, error) {
	var users []models.User

	query := `SELECT * FROM users;`

	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *usersRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM users WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *usersRepository) UpdateUser(ctx context.Context, user models.User) error {
	query := `UPDATE users SET first_name = $1, last_name = $2, department = $3, email = $4, phone = $5, user_pic = $6 WHERE user_id = $7`

	_, err := r.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Department, user.Email, user.Phone, user.Userpic, user.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *usersRepository) ListDepartmentUsers(ctx context.Context, userID string) ([]*models.User, error) {
	var depID uuid.UUID

	err := r.db.GetContext(ctx, &depID, "SELECT department FROM users WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	query := `
						SELECT 
            u.user_id,
            u.first_name,
            u.last_name,
            u.department,
            u.email,
            u.phone,
            u.user_pic,
            case 
                when u.latest_ticket is null then '{}'::jsonb
                else jsonb_build_object(
                    'status', t.status,
                    'workstarted_at', t.workstarted_at,
                    'workfinished_at', t.workfinished_at,
                    'client_name', c.title
                )
            end as properties,
            (
                select count(*) 
                from tickets 
                where executor = u.user_id 
                and status not in ('closed', 'cancelled')
            ) as active_tickets
        from 
            users u
        left join tickets t on u.latest_ticket = t.id
        left join clients c on t.client = c.id
        where 
            u.department = $1
	`
	var users []*models.User

	err = r.db.SelectContext(ctx, &users, query, depID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
