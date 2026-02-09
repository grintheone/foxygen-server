package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Manufacturer struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Title string    `json:"title" db:"title"`
}

type ManufacturerRepo interface {
	AddNewManufacturer(ctx context.Context, manufacturer Manufacturer) error
}

type manufacturerRepo struct {
	db *sqlx.DB
}

func NewManufacturerRepo(db *sqlx.DB) *manufacturerRepo {
	return &manufacturerRepo{db}
}

func (r *manufacturerRepo) AddNewManufacturer(ctx context.Context, manufacturer Manufacturer) error {
	query := `INSERT INTO manufacturers (id, title) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, manufacturer.ID, manufacturer.Title)
	if err != nil {
		return err
	}

	return nil
}
