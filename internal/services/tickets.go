package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type TicketService struct {
	repo repository.TicketsRepository
}

func NewTicketService(r repository.TicketsRepository) *TicketService {
	return &TicketService{repo: r}
}

func (s *TicketService) ListAllTickets(ctx context.Context, executorID string) ([]*models.TicketCard, error) {
	tickets, err := s.repo.ListAllTickets(ctx, executorID)
	if err != nil {
		return nil, fmt.Errorf("service error getting all tickets: %w", err)
	}

	return tickets, nil
}

func (s *TicketService) GetTicketByID(ctx context.Context, uuid uuid.UUID) (*models.TicketSinglePage, error) {
	ticket, err := s.repo.GetTicketByID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service error getting ticket by ID: %w", err)
	}

	return ticket, nil
}

func (s *TicketService) DeleteTicketByID(ctx context.Context, uuid uuid.UUID) error {
	err := s.repo.DeleteTicketByID(ctx, uuid)
	if err != nil {
		return fmt.Errorf("service error deleting ticket by ID: %w", err)
	}

	return nil
}

func (s *TicketService) CreateNewTicket(ctx context.Context, payload models.RawTicket) (*models.RawTicket, error) {
	created, err := s.repo.CreateNewTicket(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("service error creating new ticket: %w", err)
	}

	return created, nil
}

func (s *TicketService) UpdateTicketInfo(ctx context.Context, uuid uuid.UUID, payload models.TicketUpdates) (*models.TicketSinglePage, error) {
	updated, err := s.repo.UpdateTicketInfo(ctx, uuid, payload)
	if err != nil {
		return nil, fmt.Errorf("service error updating ticket info: %w", err)
	}

	return updated, nil
}

func (s *TicketService) GetReasonInfoByID(ctx context.Context, id string) (*models.TicketReason, error) {
	reasonInfo, err := s.repo.GetReasonInfoByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service error getting reason info: %w", err)
	}

	return reasonInfo, nil
}

func (s *TicketService) GetTicketContactPerson(ctx context.Context, uuid uuid.UUID) (*models.Contact, error) {
	contact, err := s.repo.GetTicketContactPerson(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service error getting ticket contact: %w", err)
	}

	return contact, nil
}

func (s *TicketService) GetTicketsByField(ctx context.Context, field string, fieldUUID uuid.UUID) ([]*models.RawTicket, error) {
	tickets, err := s.repo.GetTicketsByField(ctx, field, fieldUUID)
	if err != nil {
		return nil, fmt.Errorf("service error getting client tickets: %w", err)
	}

	return tickets, nil
}
