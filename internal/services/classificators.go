package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type ClassificatorService struct {
	repo repository.ClassificatorsRepository
}

func NewClassificatorService(r repository.ClassificatorsRepository) *ClassificatorService {
	return &ClassificatorService{repo: r}
}

func (s *ClassificatorService) GetClassificatorByID(ctx context.Context, uuid uuid.UUID) (*models.Classificator, error) {
	classificator, err := s.repo.GetClassificatorByID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service error getting classificator: %w", err)
	}

	return classificator, nil
}

func (s *ClassificatorService) GetDevicesByClassificatorID(ctx context.Context, uuid uuid.UUID) (*[]models.Device, error) {
	devices, err := s.repo.GetDevicesByClassificatorID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service error gettting classificator devices: %w", err)
	}

	return devices, nil
}

func (s *ClassificatorService) NewClassificator(ctx context.Context, payload models.Classificator) (*models.Classificator, error) {
	created, err := s.repo.NewClassificator(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("service error creating new classificator: %w", err)
	}

	return created, nil
}

func (s *ClassificatorService) RemoveClassificatorByID(ctx context.Context, uuid uuid.UUID) error {
	err := s.repo.RemoveClassificatorByID(ctx, uuid)
	if err != nil {
		return fmt.Errorf("service error deliting classificator: %w", err)
	}

	return nil
}

func (s *ClassificatorService) UpdateClassificatorInfo(ctx context.Context, uuid uuid.UUID, payload models.ClassificatorUpdate) (*models.Classificator, error) {
	updated, err := s.repo.UpdateClassificatorInfo(ctx, uuid, payload)
	if err != nil {
		return nil, fmt.Errorf("service error updating classificator: %w", err)
	}

	return updated, nil
}
