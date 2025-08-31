package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type DevicesRepository interface {
	GetAllDevices(ctx context.Context) (*[]models.Device, error)
	GetDeviceByID(ctx context.Context, uuid uuid.UUID) (*models.Device, error)
	RemoveDeviceByID(ctx context.Context, uuid uuid.UUID) error
	CreateNewDevice(ctx context.Context, payload models.Device) (*models.Device, error)
	UpdateDeviceByID(ctx context.Context, uuid uuid.UUID, payload models.DeviceUpdates) (*models.Device, error)
}

type deviceRepository struct {
	db *sqlx.DB
}

func NewDeviceRepository(db *sqlx.DB) *deviceRepository {
	return &deviceRepository{db}
}

func (r *deviceRepository) GetAllDevices(ctx context.Context) (*[]models.Device, error) {
	query := `SELECT * FROM devices`

	var devices []models.Device
	err := r.db.SelectContext(ctx, &devices, query)
	if err != nil {
		return nil, err
	}

	return &devices, nil
}

func (r *deviceRepository) GetDeviceByID(ctx context.Context, uuid uuid.UUID) (*models.Device, error) {
	var device models.Device
	err := r.db.GetContext(ctx, &device, `SELECT * FROM devices WHERE id = $1`, uuid)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *deviceRepository) RemoveDeviceByID(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM devices WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (r *deviceRepository) CreateNewDevice(ctx context.Context, payload models.Device) (*models.Device, error) {
	query := `
		INSERT INTO devices (classificator, serial_number, properties, connected_to_lis, is_used)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`

	var device models.Device

	err := r.db.GetContext(ctx, &device, query, payload.Classificator, payload.SerialNumber, payload.Properties, payload.ConntectedToLIS, payload.IsUsed)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *deviceRepository) UpdateDeviceByID(ctx context.Context, uuid uuid.UUID, payload models.DeviceUpdates) (*models.Device, error) {
	existing, err := r.GetDeviceByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if payload.Classificator != nil {
		existing.Classificator = *payload.Classificator
	}
	if payload.SerialNumber != nil {
		existing.SerialNumber = *payload.SerialNumber
	}
	if payload.Properties != nil {
		existing.Properties = *payload.Properties
	}
	if payload.ConntectedToLIS != nil {
		existing.ConntectedToLIS = *payload.ConntectedToLIS
	}
	if payload.IsUsed != nil {
		existing.IsUsed = *payload.IsUsed
	}

	query := `
		UPDATE devices
		SET classificator = :classificator, serial_number = :serial_number, properties = :properties, connected_to_lis = :connected_to_lis, is_used = :is_used
		WHERE id = :id
	`

	_, err = r.db.NamedExecContext(ctx, query, &existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}
