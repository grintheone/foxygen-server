package services

import (
	"context"
	"fmt"
	"strconv"

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

func (s *AgreementService) GetAgreementsByField(ctx context.Context, uuid uuid.UUID, field string, active string) ([]*models.AgreementCard, error) {
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

	isActive, err := strconv.ParseBool(active)
	if err != nil {
		return nil, fmt.Errorf("active flag is not a boolean %w", err)
	}

	agreements, err := s.repo.GetAgreementsByField(ctx, field, uuid, isActive)
	if err != nil {
		return nil, fmt.Errorf("service error getting client agreements: %w", err)
	}

	return agreements, nil
}
