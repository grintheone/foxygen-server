package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type ClassificatorsRepository interface {
	ListClassificators(ctx context.Context) (*[]models.Classificator, error)
	NewClassificator(ctx context.Context, payload models.Classificator) error
	GetClassificatorByID(ctx context.Context, uuid uuid.UUID) (*models.Classificator, error)
	GetDevicesByClassificatorID(ctx context.Context, uuid uuid.UUID) (*[]models.Device, error)
	RemoveClassificatorByID(ctx context.Context, uuid uuid.UUID) error
	UpdateClassificatorInfo(ctx context.Context, uuid uuid.UUID, payload models.ClassificatorUpdate) (*models.Classificator, error)
}

type classificatorRepository struct {
	db *sqlx.DB
}

func NewClassificatorRepository(db *sqlx.DB) *classificatorRepository {
	return &classificatorRepository{db}
}

func (r *classificatorRepository) ListClassificators(ctx context.Context) (*[]models.Classificator, error) {
	query := `SELECT * FROM classificators`

	var classificators []models.Classificator
	err := r.db.SelectContext(ctx, &classificators, query)
	if err != nil {
		return nil, err
	}

	return &classificators, nil
}

func (r *classificatorRepository) GetClassificatorByID(ctx context.Context, uuid uuid.UUID) (*models.Classificator, error) {
	query := `SELECT * FROM classificators WHERE id = $1`

	var c models.Classificator

	err := r.db.GetContext(ctx, &c, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

func (r *classificatorRepository) GetDevicesByClassificatorID(ctx context.Context, uuid uuid.UUID) (*[]models.Device, error) {
	query := `SELECT * FROM devices WHERE classificators = $1`

	var devices []models.Device
	err := r.db.SelectContext(ctx, &devices, query, uuid)
	if err != nil {
		return nil, err
	}

	return &devices, nil
}

func (r *classificatorRepository) NewClassificator(ctx context.Context, payload models.Classificator) error {
	query := `
		INSERT INTO classificators (id, title, manufacturer, research_type, registration_certificate, maintenance_regulations, attachments, images)
    VALUES (:id, :title, :manufacturer, :research_type, :registration_certificate, :maintenance_regulations, :attachments, :images)
	`

	_, err := r.db.NamedExecContext(ctx, query, payload)
	if err != nil {
		return err
	}

	return nil
}

func (r *classificatorRepository) RemoveClassificatorByID(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM classificators WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (r *classificatorRepository) UpdateClassificatorInfo(ctx context.Context, uuid uuid.UUID, payload models.ClassificatorUpdate) (*models.Classificator, error) {
	existing, err := r.GetClassificatorByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if payload.Title != nil {
		existing.Title = *payload.Title
	}
	if payload.Manufacturer != nil {
		existing.Manufacturer = payload.Manufacturer
	}
	if payload.ResearchType != nil {
		existing.ResearchType = payload.ResearchType
	}
	if payload.RegistrationCertificate != nil {
		existing.RegistrationCertificate = *payload.RegistrationCertificate
	}
	if payload.MaintenanceRegulations != nil {
		existing.MaintenanceRegulations = payload.MaintenanceRegulations
	}
	if payload.Attachments != nil {
		existing.Attachments = *payload.Attachments
	}
	if payload.Images != nil {
		existing.Images = *payload.Images
	}

	query := `
		UPDATE classificators
		SET title = :title, manufacturer = :manufacturer, research_type = :research_type, registration_certificate = :registration_certificate, maintenance_regulations = :maintenance_regulations, attachments = :attachments, images = :images
		WHERE id = :id
	`

	_, err = r.db.NamedExecContext(ctx, query, &existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}
