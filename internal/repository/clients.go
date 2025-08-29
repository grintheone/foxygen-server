package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type ClientsRepository interface {
	ListClients(ctx context.Context) (*[]models.Client, error)
	CreateClient(ctx context.Context, payload models.Client) (*models.Client, error)
	UpdateClient(ctx context.Context, uuid uuid.UUID, payload models.ClientUpdate) (*models.Client, error)
	DeleteClient(ctx context.Context, uuid uuid.UUID) error
}

type clientsRepository struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) ClientsRepository {
	return &clientsRepository{db}
}

func (r *clientsRepository) ListClients(ctx context.Context) (*[]models.Client, error) {
	var clients []models.Client

	query := `SELECT * FROM clients;`

	err := r.db.SelectContext(ctx, &clients, query)
	if err != nil {
		return nil, err
	}

	return &clients, nil
}

func (r *clientsRepository) CreateClient(ctx context.Context, payload models.Client) (*models.Client, error) {
	query := `
        INSERT INTO clients (title, region, address, location, laboratory_system, manager)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING *
    `

	var client models.Client

	err := r.db.GetContext(ctx, &client, query, payload.Title, payload.Region, payload.Address, payload.Location, payload.LaboratorySystem, payload.Manager)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *clientsRepository) UpdateClient(ctx context.Context, uuid uuid.UUID, payload models.ClientUpdate) (*models.Client, error) {
	var existing models.Client

	err := r.db.GetContext(ctx, &existing, `SELECT * FROM clients WHERE id = $1`, uuid)
	if err != nil {
		return nil, err
	}

	if payload.Title != nil {
		existing.Title = *payload.Title
	}
	if payload.Region != nil {
		existing.Region = *payload.Region
	}
	if payload.Address != nil {
		existing.Address = *payload.Address
	}
	if payload.Location != nil {
		existing.Location = *payload.Location
	}
	if payload.LaboratorySystem != nil {
		existing.LaboratorySystem = *payload.LaboratorySystem
	}
	if payload.Manager != nil {
		existing.Manager = *payload.Manager
	}

	query := `
        UPDATE clients
        SET title = :title, region = :region, address = :address, location = :location, laboratory_system = :laboratory_system
        WHERE id = :id
    `

	_, err = r.db.NamedExecContext(ctx, query, &existing)
	if err != nil {
		return nil, err
	}

	return &existing, nil
}

func (r *clientsRepository) DeleteClient(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM clients WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	return nil
}
