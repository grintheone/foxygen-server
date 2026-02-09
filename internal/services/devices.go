package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type DeviceService struct {
	repo repository.DevicesRepository
}

func NewDeviceService(r repository.DevicesRepository) *DeviceService {
	return &DeviceService{repo: r}
}

func (s *DeviceService) GetAllDevices(ctx context.Context) (*[]models.Device, error) {
	devices, err := s.repo.GetAllDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("service error getting all devices: %w", err)
	}
	return devices, nil
}

func (s *DeviceService) GetDeviceByID(ctx context.Context, uuid uuid.UUID) (*models.DeviceSinglePage, error) {
	device, err := s.repo.GetDeviceByID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service error getting device by ID: %w", err)
	}

	return device, nil
}

func (s *DeviceService) RemoveDeviceByID(ctx context.Context, uuid uuid.UUID) error {
	err := s.repo.RemoveDeviceByID(ctx, uuid)
	if err != nil {
		return fmt.Errorf("service error deleting device: %w", err)
	}

	return nil
}

func (s *DeviceService) CreateNewDevice(ctx context.Context, payload models.Device) error {
	err := s.repo.CreateNewDevice(ctx, payload)
	if err != nil {
		return fmt.Errorf("service error creating device: %w", err)
	}

	return nil
}

func (s *DeviceService) UpdateDeviceByID(ctx context.Context, uuid uuid.UUID, payload models.DeviceUpdates) (*models.DeviceSinglePage, error) {
	updated, err := s.repo.UpdateDeviceByID(ctx, uuid, payload)
	if err != nil {
		return nil, fmt.Errorf("service error updating device: %w", err)
	}

	return updated, nil
}

func (s *DeviceService) GetDeviceRemoteOptions(ctx context.Context, uuid uuid.UUID) ([]*models.DeviceRemoteOption, error) {
	options, err := s.repo.GetDeviceRemoteOptions(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service error getting remote options: %w", err)
	}

	return options, nil
}
