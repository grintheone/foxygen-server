package models

import (
	"github.com/google/uuid"
)

type Account struct {
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Username     string    `json:"username" db:"username"`
	Database     string    `json:"database" db:"database"`
	Disabled     bool      `json:"disabled" db:"disabled"`
	PasswordHash string    `json:"-" db:"password_hash"` // The `-` hides this field in JSON responses
	Roles        []Role    `json:"roles,omitempty"`      // This is for convenience when fetching a user with their roles
}

// Role represents a user role
type Role struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

// AccountRole represents the junction table. Used for inserting associations.
type AccountRole struct {
	UserID uuid.UUID `db:"user_id"`
	RoleID int       `db:"role_id"`
}

// GetRoleNames returns a slice of role names for the account.
func (a *Account) GetRoleNames() []string {
	var roleNames []string
	for _, role := range a.Roles {
		roleNames = append(roleNames, role.Name)
	}
	return roleNames
}
