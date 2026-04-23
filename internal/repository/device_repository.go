package repository

import (
	"context"
	"errors"
	"fmt"

	"dev.paultergust/devices-api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type DeviceRepository struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(ctx context.Context, d *model.Device) error {
	query := `
        INSERT INTO devices (id, name, brand, state, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := r.db.Exec(ctx, query,
		d.ID, d.Name, d.Brand, d.State, d.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert device: %w", err)
	}
	return nil
}

func (r *DeviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Device, error) {
	query := `
        SELECT id, name, brand, state, created_at
        FROM devices
        WHERE id = $1
    `

	var d model.Device
	err := r.db.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.Name, &d.Brand, &d.State, &d.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get device by id: %w", err)
	}

	return &d, nil
}

func (r *DeviceRepository) List(ctx context.Context, brand, state *string) ([]model.Device, error) {
	query := "SELECT id, name, brand, state, created_at FROM devices WHERE 1=1"
	args := []interface{}{}
	i := 1

	if brand != nil {
		query += fmt.Sprintf(" AND brand = $%d", i)
		args = append(args, *brand)
		i++
	}

	if state != nil {
		query += fmt.Sprintf(" AND state = $%d", i)
		args = append(args, *state)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list devices query: %w", err)
	}
	defer rows.Close()

	var devices []model.Device

	for rows.Next() {
		var d model.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.Brand, &d.State, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan device: %w", err)
		}
		devices = append(devices, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate devices: %w", err)
	}

	return devices, nil
}

func (r *DeviceRepository) Update(ctx context.Context, d *model.Device) error {
	query := `
        UPDATE devices
        SET name=$1, brand=$2, state=$3
        WHERE id=$4
    `
	cmdTag, err := r.db.Exec(ctx, query,
		d.Name, d.Brand, d.State, d.ID,
	)
	if err != nil {
		return fmt.Errorf("update device: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *DeviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmdTag, err := r.db.Exec(ctx, "DELETE FROM devices WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("delete device: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
