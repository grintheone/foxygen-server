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
