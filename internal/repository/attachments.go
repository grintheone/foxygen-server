package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *models.Attachment) error
	CreateBulk(ctx context.Context, attachments []*models.Attachment) error
	GetAttachmentsByRefID(ctx context.Context, refID uuid.UUID) ([]*models.Attachment, error)
	GetAttachmentByID(ctx context.Context, id string) (*models.Attachment, error)
	// Delete(ctx context.Context, id int) error
}

type attachmentRepository struct {
	db *sqlx.DB
}

func NewAttachmentRepository(db *sqlx.DB) AttachmentRepository {
	return &attachmentRepository{db}
}

func (r *attachmentRepository) Create(ctx context.Context, attachment *models.Attachment) error {
	query := `
		INSERT INTO attachments (id, name, media_type, ext, ref_id)
		VALUES (:id, :name, :media_type, :ext, :ref_id)
	`

	_, err := r.db.NamedExecContext(ctx, query, attachment)
	if err != nil {
		return err
	}

	return nil
}

func (r *attachmentRepository) CreateBulk(ctx context.Context, attachments []*models.Attachment) error {
	if len(attachments) == 0 {
		return nil
	}

	query := `
			INSERT INTO attachments (id, name, media_type, ext, ref_id)
			VALUES (:id, :name, :media_type, :ext, :ref_id)
			`

	_, err := r.db.NamedExecContext(ctx, query, attachments)
	return err
}

func (r *attachmentRepository) GetAttachmentsByRefID(ctx context.Context, refID uuid.UUID) ([]*models.Attachment, error) {
	query := `
		SELECT * FROM attachments WHERE ref_id = $1
	`

	var attachments []*models.Attachment

	err := r.db.SelectContext(ctx, &attachments, query, refID)
	if err != nil {
		return nil, err
	}

	return attachments, nil
}

func (r *attachmentRepository) GetAttachmentByID(ctx context.Context, id string) (*models.Attachment, error) {
	query := `SELECT * FROM attachments WHERE id = $1`

	var attachment models.Attachment

	err := r.db.GetContext(ctx, &attachment, query, id)
	if err != nil {
		return nil, err
	}

	return &attachment, nil
}
