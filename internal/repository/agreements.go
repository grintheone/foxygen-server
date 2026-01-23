package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type AgreementRepo interface {
	GetAgreementsByField(ctx context.Context, field string, uuid uuid.UUID, active bool) ([]*models.AgreementCard, error)
}

type agreementRepo struct {
	db *sqlx.DB
}

func NewAgreementRepo(db *sqlx.DB) AgreementRepo {
	return &agreementRepo{db}
}

func (r *agreementRepo) GetAgreementsByField(ctx context.Context, field string, uuid uuid.UUID, active bool) ([]*models.AgreementCard, error) {
	var agreements []*models.AgreementCard
	query := fmt.Sprintf(`
		SELECT
	 		a.*,
			c.title as client_name,
			c.address as client_address,
			d.serial_number as device_serial,
			cl.title as device_name
		FROM agreements a
		LEFT JOIN clients c ON a.actual_client = c.id
		LEFT JOIN devices d ON a.device = d.id
		LEFT JOIN classificators cl ON d.classificator = cl.id
		WHERE %s = $1 AND is_active = $2
	`, field)

	err := r.db.SelectContext(ctx, &agreements, query, uuid, active)
	if err != nil {
		return nil, err
	}

	return agreements, nil
}
