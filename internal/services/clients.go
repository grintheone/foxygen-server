package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type ClientService struct {
	repo repository.ClientsRepository
}

func NewClientService(r repository.ClientsRepository) *ClientService {
	return &ClientService{repo: r}
}

func (s *ClientService) ListClients(ctx context.Context, limit int, offset int, sortByTitle bool) (*[]models.Client, error) {
	clients, err := s.repo.ListClients(ctx, limit, offset, sortByTitle)
	if err != nil {
		return nil, fmt.Errorf("service error fetching all clients: %w", err)
	}
	return clients, nil
}

func (s *ClientService) CreateClient(ctx context.Context, payload models.Client) error {
	err := s.repo.CreateClient(ctx, payload)
	if err != nil {
		return fmt.Errorf("service error creating a client: %w", err)
	}

	return nil
}

func (s *ClientService) UpdateClient(ctx context.Context, uuid uuid.UUID, payload models.ClientUpdate) (*models.Client, error) {
	client, err := s.repo.UpdateClient(ctx, uuid, payload)
	if err != nil {
		return nil, fmt.Errorf("service error updating a client: %w", err)
	}

	return client, nil
}

func (s *ClientService) DeleteClient(ctx context.Context, uuid uuid.UUID) error {
	err := s.repo.DeleteClient(ctx, uuid)
	if err != nil {
		return fmt.Errorf("service error deleting a client: %w", err)
	}

	return nil
}

func (s *ClientService) GetClientByID(ctx context.Context, uuid uuid.UUID) (*models.Client, error) {
	client, err := s.repo.GetClientByID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service error getting client by ID: %w", err)
	}

	return client, nil
}
