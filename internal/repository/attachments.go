package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *models.Attachment) error
	GetAttachmentsByRefID(ctx context.Context, refID uuid.UUID) (*[]models.Attachment, error)
	GetAttachmentByID(ctx context.Context, id int) (*models.Attachment, error)
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
		INSERT INTO attachments (file_name, original_name, file_size, file_path, mime_type, ref_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query, attachment.FileName, attachment.OriginalName, attachment.FileSize, attachment.FilePath, attachment.MimeType, attachment.RefID)
	if err != nil {
		return err
	}

	return nil
}

func (r *attachmentRepository) GetAttachmentsByRefID(ctx context.Context, refID uuid.UUID) (*[]models.Attachment, error) {
	query := `
		SELECT * FROM attachments WHERE ref_id = $1
	`

	var attachments []models.Attachment

	err := r.db.SelectContext(ctx, &attachments, query, refID)
	if err != nil {
		return nil, err
	}

	return &attachments, nil
}

func (r *attachmentRepository) GetAttachmentByID(ctx context.Context, id int) (*models.Attachment, error) {
	query := `SELECT * FROM attachments WHERE id = $1`

	var attachment models.Attachment

	err := r.db.GetContext(ctx, &attachment, query, id)
	if err != nil {
		return nil, err
	}

	return &attachment, nil
}
