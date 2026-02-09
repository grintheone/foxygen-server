package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ResearchType struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Title string    `json:"title" db:"title"`
}

type ResearchTypeRepo interface {
	AddNewResearchType(ctx context.Context, researchType ResearchType) error
}

type researchTypeRepo struct {
	db *sqlx.DB
}

func NewResearchTypeRepo(db *sqlx.DB) *researchTypeRepo {
	return &researchTypeRepo{db}
}

func (r *researchTypeRepo) AddNewResearchType(ctx context.Context, researchType ResearchType) error {
	query := `INSERT INTO research_type (id, title) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, researchType.ID, researchType.Title)
	if err != nil {
		return err
	}

	return nil
}
