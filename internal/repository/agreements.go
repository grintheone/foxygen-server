package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type AgreementRepo interface {
	GetClientAgreements(ctx context.Context, clientID uuid.UUID) ([]*models.Agreement, error)
}

type agreementRepo struct {
	db *sqlx.DB
}

func NewAgreementRepo(db *sqlx.DB) AgreementRepo {
	return &agreementRepo{db}
}

func (r *agreementRepo) GetClientAgreements(ctx context.Context, clientID uuid.UUID) ([]*models.Agreement, error) {
	query := `SELECT * FROM agreements WHERE actual_client = $1`
	var agreements []*models.Agreement

	err := r.db.SelectContext(ctx, &agreements, query, clientID)
	if err != nil {
		return nil, err
	}

	return agreements, nil
}
