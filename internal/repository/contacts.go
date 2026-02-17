package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type ContactsRepository interface {
	GetAllByClientID(ctx context.Context, uuid uuid.UUID) (*[]models.ContactRow, error)
	CreateContact(ctx context.Context, data models.Contact) error
	DeleteContact(ctx context.Context, uuid uuid.UUID) error
	GetContactByID(ctx context.Context, uuid uuid.UUID) (*models.Contact, error)
	UpdateContact(ctx context.Context, uuid uuid.UUID, payload models.ContactUpdate) (*models.Contact, error)
}

type contactRepository struct {
	db *sqlx.DB
}

func NewContactRepository(db *sqlx.DB) *contactRepository {
	return &contactRepository{db}
}

func (r *contactRepository) GetAllByClientID(ctx context.Context, uuid uuid.UUID) (*[]models.ContactRow, error) {
	query := `
		SELECT id, name, position, phone, email
		FROM contacts
		WHERE client_id = $1
	`
	var contacts []models.ContactRow
	err := r.db.SelectContext(ctx, &contacts, query, uuid)
	if err != nil {
		return nil, err
	}

	return &contacts, nil
}

func (r *contactRepository) GetContactByID(ctx context.Context, uuid uuid.UUID) (*models.Contact, error) {
	query := `SELECT FROM contacts WHERE id = $1`

	var contact models.Contact

	err := r.db.GetContext(ctx, &contact, query, uuid)
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (r *contactRepository) CreateContact(ctx context.Context, data models.Contact) error {
	query := `
		INSERT INTO contacts (id, name, position, phone, email, client_id)
		VALUES (:id, :name, :position, :phone, :email, :client_id)
	`
	_, err := r.db.NamedExecContext(ctx, query, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *contactRepository) DeleteContact(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM contacts WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (r *contactRepository) UpdateContact(ctx context.Context, uuid uuid.UUID, payload models.ContactUpdate) (*models.Contact, error) {
	// Get existing contact to merge changes
	var existing models.Contact
	err := r.db.GetContext(ctx, &existing, `SELECT * FROM contacts WHERE id = $1`, uuid)
	if err != nil {
		return nil, err
	}

	if payload.Name != nil {
		existing.Name = *payload.Name
	}

	if payload.Position != nil {
		existing.Position = *payload.Position
	}

	if payload.Phone != nil {
		existing.Phone = *payload.Phone
	}

	if payload.Email != nil {
		existing.Email = *payload.Email
	}

	if payload.ClientID != nil {
		existing.ClientID = *payload.ClientID
	}

	query := `
        UPDATE contacts
        SET name = :name, position = :position, phone = :phone, email = :email, client_id = :client_id
        WHERE id = :id
    `

	_, err = r.db.NamedExecContext(ctx, query, &existing)
	if err != nil {
		return nil, err
	}

	return &existing, nil
}
