package models

import (
	"github.com/google/uuid"
)

type Account struct {
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Username     string    `json:"username" db:"username"`
	Disabled     bool      `json:"disabled" db:"disabled"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role"` // This is for convenience when fetching a user with their roles
}

// Role represents a user role
type Role struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type AccountRole struct {
	UserID uuid.UUID `db:"user_id"`
	RoleID int       `db:"role_id"`
}
