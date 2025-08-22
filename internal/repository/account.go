package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db}
}

func (r *AccountRepository) CreateAccountWithRoles(ctx context.Context, account *models.Account, roleID int) (*models.Account, error) {
	// Begin a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	// Defer a rollback in case anything fails. It's a no-op if the tx is committed.
	defer tx.Rollback()

	// Generate a new UUID for the user
	newUserID := uuid.New()
	account.UserID = newUserID

	// Set current account role
	switch roleID {
	case 1:
		account.Role = "admin"
	case 2:
		account.Role = "coordinator"
	case 3:
		account.Role = "user"
	}

	// Insert the new user into the 'accounts' table
	query := `INSERT INTO accounts (user_id, username, database, password_hash) VALUES ($1, $2, $3, $4)`
	_, err = tx.ExecContext(ctx, query, account.UserID, account.Username, account.Database, account.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("could not insert account: %w", err)
	}

	// Insert the role associations into the 'account_roles' table
	accRoleQuery := `INSERT INTO account_roles (user_id, role_id) VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, accRoleQuery, account.UserID, roleID)
	if err != nil {
		return nil, fmt.Errorf("could not assign roles to account: %w", err)
	}

	// If everything went well, commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return account, nil
}

func (r *AccountRepository) GetByUsername(ctx context.Context, username string) (*models.Account, error) {
	query := `
	   	SELECT
	        a.user_id,
	        a.username,
	        a.database,
	        a.disabled,
	        a.password_hash,
	        r.name as role
	    FROM accounts a
	    LEFT JOIN account_roles ar ON a.user_id = ar.user_id
	    LEFT JOIN roles r ON ar.role_id = r.id
	    WHERE a.username = $1;
    `

	var account models.Account

	err := r.db.GetContext(ctx, &account, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return a nil account and no error if not found.
			// This allows the service to handle "not found" as a specific case.
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get account by username: %w", err)
	}

	return &account, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.Account, error) {
	query := `
	   	SELECT
	        a.user_id,
	        a.username,
	        a.database,
	        a.disabled,
	        a.password_hash,
	        r.name as role
	    FROM accounts a
	    LEFT JOIN account_roles ar ON a.user_id = ar.user_id
	    LEFT JOIN roles r ON ar.role_id = r.id
	    WHERE a.user_id = $1;
    `

	var account models.Account

	err := r.db.GetContext(ctx, &account, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get account by ID: %w", err)
	}

	return &account, nil
}

func (r *AccountRepository) ChangeAccountPassword(ctx context.Context, userID uuid.UUID, hash string) error {
	query := `UPDATE accounts SET password_hash = $1 WHERE user_id = $2`

	_, err := r.db.ExecContext(ctx, query, hash, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *AccountRepository) ChangeAccountStatus(ctx context.Context, userID uuid.UUID, disabled bool) error {
	query := `UPDATE accounts SET disabled = $1 WHERE user_id = $2`

	_, err := r.db.ExecContext(ctx, query, disabled, userID)
	if err != nil {
		return err
	}

	return nil
}
