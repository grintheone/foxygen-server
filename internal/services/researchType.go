package services

import (
	"context"
	"fmt"

	"github.com/grintheone/foxygen-server/internal/repository"
)

type ResearchTypeService struct {
	repo repository.ResearchTypeRepo
}

func NewResearchTypeService(repo repository.ResearchTypeRepo) *ResearchTypeService {
	return &ResearchTypeService{repo}
}

func (s *ResearchTypeService) ListAllResearchTypes(ctx context.Context) ([]*repository.ResearchType, error) {
	researchTypes, err := s.repo.ListAllResearchTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("service error fetching research types: %w", err)
	}

	return researchTypes, nil
}
