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
	ListAllManufacturers(ctx context.Context) ([]*Manufacturer, error)
	AddNewManufacturer(ctx context.Context, manufacturer Manufacturer) error
}

type manufacturerRepo struct {
	db *sqlx.DB
}

func NewManufacturerRepo(db *sqlx.DB) *manufacturerRepo {
	return &manufacturerRepo{db}
}

func (r *manufacturerRepo) ListAllManufacturers(ctx context.Context) ([]*Manufacturer, error) {
	var manufacturers []*Manufacturer

	err := r.db.SelectContext(ctx, &manufacturers, `SELECT id, title FROM manufacturers ORDER BY title ASC`)
	if err != nil {
		return nil, err
	}

	return manufacturers, nil
}

func (r *manufacturerRepo) AddNewManufacturer(ctx context.Context, manufacturer Manufacturer) error {
	query := `INSERT INTO manufacturers (id, title) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, manufacturer.ID, manufacturer.Title)
	if err != nil {
		return err
	}

	return nil
}
