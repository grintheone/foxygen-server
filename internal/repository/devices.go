package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type DevicesRepository interface {
	GetAllDevices(ctx context.Context, limit int, offset int, sortByTitle bool) (*[]models.Device, error)
	GetDeviceByID(ctx context.Context, uuid uuid.UUID) (*models.DeviceSinglePage, error)
	RemoveDeviceByID(ctx context.Context, uuid uuid.UUID) error
	CreateNewDevice(ctx context.Context, payload models.Device) error
	UpdateDeviceByID(ctx context.Context, uuid uuid.UUID, payload models.DeviceUpdates) (*models.DeviceSinglePage, error)
	GetDeviceRemoteOptions(ctx context.Context, uuid uuid.UUID) ([]*models.DeviceRemoteOption, error)
}

type deviceRepository struct {
	db *sqlx.DB
}

func NewDeviceRepository(db *sqlx.DB) *deviceRepository {
	return &deviceRepository{db}
}

func (r *deviceRepository) GetAllDevices(ctx context.Context, limit int, offset int, sortByTitle bool) (*[]models.Device, error) {
	query := `SELECT * FROM devices`
	if sortByTitle {
		query += ` ORDER BY serial_number ASC`
	}
	query += ` LIMIT $1 OFFSET $2`

	var devices []models.Device
	err := r.db.SelectContext(ctx, &devices, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return &devices, nil
}

func (r *deviceRepository) GetDeviceByID(ctx context.Context, uuid uuid.UUID) (*models.DeviceSinglePage, error) {
	var device models.DeviceSinglePage
	query := `
		SELECT
	 		d.*,
			cl.title as classificator
		FROM devices d
		LEFT JOIN classificators cl ON d.classificator = cl.id
		WHERE d.id = $1
	`

	err := r.db.GetContext(ctx, &device, query, uuid)
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

func (r *deviceRepository) CreateNewDevice(ctx context.Context, payload models.Device) error {
	query := `
		INSERT INTO devices (id, classificator, serial_number, properties, connected_to_lis, is_used)
	VALUES (:id, :classificator, :serial_number, :properties, :connected_to_lis, :is_used)
	`

	_, err := r.db.NamedExecContext(ctx, query, payload)
	if err != nil {
		return err
	}

	return nil
}

func (r *deviceRepository) UpdateDeviceByID(ctx context.Context, uuid uuid.UUID, payload models.DeviceUpdates) (*models.DeviceSinglePage, error) {
	query := `
		UPDATE devices
		SET
			classificator = COALESCE($2, classificator),
			serial_number = COALESCE($3, serial_number),
			properties = COALESCE($4, properties),
			connected_to_lis = COALESCE($5, connected_to_lis),
			is_used = COALESCE($6, is_used)
		WHERE id = $1
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		uuid,
		payload.Classificator,
		payload.SerialNumber,
		payload.Properties,
		payload.ConntectedToLIS,
		payload.IsUsed,
	)
	if err != nil {
		return nil, err
	}

	return r.GetDeviceByID(ctx, uuid)
}

func (r *deviceRepository) GetDeviceRemoteOptions(ctx context.Context, uuid uuid.UUID) ([]*models.DeviceRemoteOption, error) {
	query := `
		SELECT
		ra.id,
		ro.title as title
		FROM remote_access ra
		LEFT JOIN ra_options ro ON ra.parameter_id = ro.id
		WHERE device_id = $1
		`

	var options []*models.DeviceRemoteOption
	err := r.db.SelectContext(ctx, &options, query, uuid)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return options, nil
}
