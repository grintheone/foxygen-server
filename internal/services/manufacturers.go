package services

import (
	"context"
	"fmt"

	"github.com/grintheone/foxygen-server/internal/repository"
)

type ManufacturerService struct {
	repo repository.ManufacturerRepo
}

func NewManufacturerService(repo repository.ManufacturerRepo) *ManufacturerService {
	return &ManufacturerService{repo}
}

func (s *ManufacturerService) ListAllManufacturers(ctx context.Context) ([]*repository.Manufacturer, error) {
	manufacturers, err := s.repo.ListAllManufacturers(ctx)
	if err != nil {
		return nil, fmt.Errorf("service error fetching manufacturers: %w", err)
	}

	return manufacturers, nil
}
