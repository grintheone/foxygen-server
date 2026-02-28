package repository

import (
	"context"

	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type RegionsRepo interface {
	ListAllRegions(ctx context.Context) ([]*models.Region, error)
	AddNewRegion(ctx context.Context, region models.Region) error
}

type regionsRepo struct {
	db *sqlx.DB
}

func NewRegionRepo(db *sqlx.DB) *regionsRepo {
	return &regionsRepo{db}
}

func (r *regionsRepo) ListAllRegions(ctx context.Context) ([]*models.Region, error) {
	var regions []*models.Region

	err := r.db.SelectContext(ctx, &regions, `SELECT id, title FROM regions ORDER BY title ASC`)
	if err != nil {
		return nil, err
	}

	return regions, nil
}

func (r *regionsRepo) AddNewRegion(ctx context.Context, region models.Region) error {
	query := `INSERT INTO regions (id, title) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, region.ID, region.Title)
	if err != nil {
		return err
	}

	return nil
}
