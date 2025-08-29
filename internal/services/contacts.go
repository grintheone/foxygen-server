package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type ContactService struct {
	repo repository.ContactsRepository
}

func NewContactService(r repository.ContactsRepository) *ContactService {
	return &ContactService{repo: r}
}

func (s *ContactService) GetAllByClientID(ctx context.Context, uuid uuid.UUID) (*[]models.Contact, error) {
	contacts, err := s.repo.GetAllByClientID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service error getting all contacts: %w", err)
	}
	return contacts, nil
}

func (s *ContactService) CreateContact(ctx context.Context, payload models.Contact) (*models.Contact, error) {
	createdContact, err := s.repo.CreateContact(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("service error creating new contact: %w", err)
	}

	return createdContact, nil
}

func (s *ContactService) DeleteContact(ctx context.Context, uuid uuid.UUID) error {
	err := s.repo.DeleteContact(ctx, uuid)
	if err != nil {
		return fmt.Errorf("service error deliting the contact: %w", err)
	}
	return nil
}

func (s *ContactService) UpdateContact(ctx context.Context, uuid uuid.UUID, payload models.ContactUpdate) (*models.Contact, error) {
	updated, err := s.repo.UpdateContact(ctx, uuid, payload)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("service error updating the contact: %w", err)
	}
	return updated, nil
}
