package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Region struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Title string    `json:"title" db:"title"`
}

type RegionsRepo interface {
	AddNewRegion(ctx context.Context, region Region) error
}

type regionsRepo struct {
	db *sqlx.DB
}

func NewRegionRepo(db *sqlx.DB) *regionsRepo {
	return &regionsRepo{db}
}

func (r *regionsRepo) AddNewRegion(ctx context.Context, region Region) error {
	query := `INSERT INTO regions (id, title) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, region.ID, region.Title)
	if err != nil {
		return err
	}

	return nil
}
