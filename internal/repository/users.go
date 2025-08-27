package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type UsersRepository interface {
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

func (r *usersRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `SELECT * FROM users WHERE user_id = $1`

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
