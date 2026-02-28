package services

import (
	"context"
	"fmt"

	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type RegionService struct {
	repo repository.RegionsRepo
}

func NewRegionService(repo repository.RegionsRepo) *RegionService {
	return &RegionService{repo}
}

func (s *RegionService) ListAllRegions(ctx context.Context) ([]*models.Region, error) {
	regions, err := s.repo.ListAllRegions(ctx)
	if err != nil {
		return nil, fmt.Errorf("service error fetching regions: %w", err)
	}

	return regions, nil
}
