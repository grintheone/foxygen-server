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

// CreateAccountWithRoles creates a new account and assigns the provided role IDs.
// It takes a context, the account model (without UserID), and a slice of role IDs to assign.
func (r *AccountRepository) CreateAccountWithRoles(ctx context.Context, account *models.Account, roleIDs []int) (*models.Account, error) {
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

	// 1. Insert the new user into the 'accounts' table
	query := `INSERT INTO accounts (user_id, username, database, password_hash) VALUES ($1, $2, $3, $4)`
	_, err = tx.ExecContext(ctx, query, account.UserID, account.Username, account.Database, account.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("could not insert account: %w", err)
	}

	// 2. Insert the role associations into the 'account_roles' table
	if len(roleIDs) > 0 {
		// Prepare the query for bulk insert
		accRoleQuery := `INSERT INTO account_roles (user_id, role_id) VALUES `
		var values []interface{}

		// Build the query and value slice dynamically
		for i, roleID := range roleIDs {
			if i > 0 {
				accRoleQuery += ", "
			}
			// Generate placeholders ($1, $2), ($3, $4)...
			accRoleQuery += fmt.Sprintf("($%d, $%d)", (i*2)+1, (i*2)+2)
			values = append(values, account.UserID, roleID)
		}

		_, err = tx.ExecContext(ctx, accRoleQuery, values...)
		if err != nil {
			return nil, fmt.Errorf("could not assign roles to account: %w", err)
		}
	}

	// 3. If everything went well, commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	// 4. (Optional) Fetch the complete user with their roles if needed
	// You can return the account struct as is, or write a separate function to get the user by ID with roles.
	return account, nil
}

func (r *AccountRepository) GetByUsername(ctx context.Context, username string) (*models.Account, error) {
	// Base query to get the user account
	query := `
        SELECT * FROM accounts WHERE username = $1
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

	// Now fetch the roles for this user
	rolesQuery := `
        SELECT r.id, r.name, r.description
        FROM roles r
        INNER JOIN account_roles ar ON r.id = ar.role_id
        WHERE ar.user_id = $1
        ORDER BY r.id
    `

	var roles []models.Role
	err = r.db.SelectContext(ctx, &roles, rolesQuery, account.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for account: %w", err)
	}

	// Assign the fetched roles to the account
	account.Roles = roles

	return &account, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.Account, error) {
	query := `
        SELECT * FROM accounts WHERE user_id = $1
    `

	var account models.Account
	err := r.db.GetContext(ctx, &account, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get account by ID: %w", err)
	}

	// Fetch roles for this user
	rolesQuery := `
           SELECT r.id, r.name, r.description
           FROM roles r
           INNER JOIN account_roles ar ON r.id = ar.role_id
           WHERE ar.user_id = $1
           ORDER BY r.id
    `
	var roles []models.Role
	err = r.db.SelectContext(ctx, &roles, rolesQuery, account.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for account: %w", err)
	}

	account.Roles = roles

	return &account, nil
}
