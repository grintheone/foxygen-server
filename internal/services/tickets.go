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

func (s *TicketService) ListAllTickets(ctx context.Context, currentUserID string, role string) ([]*models.TicketCard, error) {
	if role == "user" {
		tickets, err := s.repo.ListAllTickets(ctx, currentUserID)
		if err != nil {
			return nil, fmt.Errorf("service error getting all tickets: %w", err)
		}

		return tickets, nil
	}

	if role == "coordinator" {
		tickets, err := s.repo.ListAllDepartmentTickets(ctx, currentUserID)
		if err != nil {
			return nil, fmt.Errorf("service error getting all tickets: %w", err)
		}

		return tickets, nil
	}

	return nil, nil
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

func (s *TicketService) CreateNewTicket(ctx context.Context, payload models.RawTicket) (*string, error) {
	created, err := s.repo.CreateNewTicket(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("service error creating new ticket: %w", err)
	}

	return created, nil
}

func (s *TicketService) GetTicketReasons(ctx context.Context) ([]*models.TicketReason, error) {
	reasons, err := s.repo.GetTicketReasons(ctx)
	if err != nil {
		return nil, fmt.Errorf("service error getting ticket reasons: %w", err)
	}

	return reasons, nil
}

func (s *TicketService) UpdateTicketInfo(ctx context.Context, payload models.TicketUpdates, userID string) error {
	err := s.repo.UpdateTicketInfo(ctx, payload, userID)
	if err != nil {
		return fmt.Errorf("service error updating ticket info: %w", err)
	}

	return nil
}

func (s *TicketService) CloseTicket(ctx context.Context, ticketInfo models.CloseTicket, currentUserID uuid.UUID) error {
	err := s.repo.CloseTicket(ctx, ticketInfo, currentUserID)
	if err != nil {
		return fmt.Errorf("service error closing ticket: %w", err)
	}

	return nil
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

func groupTicketsByMonth(tickets []*models.TicketCard, status string) map[string][]*models.TicketCard {
	var monthKey string
	grouped := make(map[string][]*models.TicketCard)
	for _, ticket := range tickets {
		if status == "closed" {
			monthKey = ticket.WorkFinishedAt.Format("2006-01")
		} else {
			monthKey = ticket.CreatedAt.Format("2006-01")
		}
		grouped[monthKey] = append(grouped[monthKey], ticket)
	}

	return grouped
}

func groupTicketsByReason(tickets []*models.TicketCard) map[string][]*models.TicketCard {
	var reasonKey string
	grouped := make(map[string][]*models.TicketCard)
	for _, ticket := range tickets {
		reasonKey = ticket.Reason
		grouped[reasonKey] = append(grouped[reasonKey], ticket)
	}

	return grouped
}

func (s *TicketService) GetTicketsByField(ctx context.Context, field string, fieldUUID uuid.UUID, filters models.TicketFilters, userID string) (*models.TicketArchiveResponseGrouped, error) {
	response, err := s.repo.GetTicketsByField(ctx, field, fieldUUID, filters, userID)
	if err != nil {
		return nil, fmt.Errorf("service error getting client tickets: %w", err)
	}

	var groupedResponse models.TicketArchiveResponseGrouped
	groupedResponse.Filters = response.Filters

	if filters.GroupBy != nil && *filters.GroupBy == "month" {
		grouped := groupTicketsByMonth(response.Tickets, filters.Status)
		groupedResponse.Tickets = grouped

		return &groupedResponse, nil
	}

	if filters.GroupBy != nil && *filters.GroupBy == "reason" {
		grouped := groupTicketsByReason(response.Tickets)
		groupedResponse.Tickets = grouped

		return &groupedResponse, nil
	}

	return &groupedResponse, nil
}
