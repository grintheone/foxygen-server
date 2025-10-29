package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type AgreementService struct {
	repo repository.AgreementRepo
}

func NewAgreementService(repo repository.AgreementRepo) *AgreementService {
	return &AgreementService{repo}
}

func (s *AgreementService) GetClientAgreements(ctx context.Context, clientID uuid.UUID) ([]*models.Agreement, error) {
	agreements, err := s.repo.GetClientAgreements(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("service error getting client agreements: %w", err)
	}

	return agreements, nil
}

func (s *AgreementService) GetAgreementsByField(ctx context.Context, field string, uuid uuid.UUID) ([]*models.AgreementCard, error) {
	allowedFields := map[string]bool{
		"client":      true,
		"device":      true,
		"distributor": true,
	}

	if !allowedFields[field] {
		return nil, fmt.Errorf("client error, field is not present in allowed fields: %s", field)
	}

	if field == "client" {
		field = "actual_client"
	}

	agreements, err := s.repo.GetAgreementsByField(ctx, field, uuid)
	if err != nil {
		return nil, fmt.Errorf("service error getting client agreements: %w", err)
	}

	return agreements, nil
}
