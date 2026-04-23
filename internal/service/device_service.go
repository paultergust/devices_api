package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dev.paultergust/devices-api/internal/dto"
	"dev.paultergust/devices-api/internal/model"

	"github.com/google/uuid"
)

var (
	ErrNotFound       = errors.New("device not found")
	ErrInvalidUpdate  = errors.New("invalid device update")
	ErrInvalidDelete  = errors.New("invalid device delete")
	ErrInvalidState   = errors.New("invalid device state")
)

type deviceRepository interface {
	Create(ctx context.Context, d *model.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Device, error)
	List(ctx context.Context, brand, state *string) ([]model.Device, error)
	Update(ctx context.Context, d *model.Device) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type DeviceService struct {
	repo deviceRepository
}

func NewDeviceService(r deviceRepository) *DeviceService {
	return &DeviceService{repo: r}
}

func (s *DeviceService) Create(ctx context.Context, req dto.CreateDeviceRequest) (*model.Device, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidState, err)
	}

	d := &model.Device{
		ID:        uuid.New(),
		Name:      req.Name,
		Brand:     req.Brand,
		State:     model.DeviceState(req.State),
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, d); err != nil {
		return nil, fmt.Errorf("create device: %w", err)
	}

	return d, nil
}

func (s *DeviceService) Get(ctx context.Context, id uuid.UUID) (*model.Device, error) {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	return d, nil
}

func (s *DeviceService) List(ctx context.Context, brand, state *string) ([]model.Device, error) {
	if state != nil {
		if err := validateState(*state); err != nil {
			return nil, err
		}
	}

	devices, err := s.repo.List(ctx, brand, state)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}

	return devices, nil
}

func (s *DeviceService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateDeviceRequest) (*model.Device, error) {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidUpdate, err)
	}

	// business rule
	if d.State == model.StateInUse {
		if req.Name != nil || req.Brand != nil {
			return nil, fmt.Errorf("%w: cannot update name or brand when device is in use", ErrInvalidUpdate)
		}
	}

	if req.Name != nil {
		d.Name = *req.Name
	}
	if req.Brand != nil {
		d.Brand = *req.Brand
	}
	if req.State != nil {
		d.State = model.DeviceState(*req.State)
	}

	if err := s.repo.Update(ctx, d); err != nil {
		return nil, fmt.Errorf("update device: %w", err)
	}

	return d, nil
}

func (s *DeviceService) Delete(ctx context.Context, id uuid.UUID) error {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNotFound, err)
	}

	// business rule
	if d.State == model.StateInUse {
		return fmt.Errorf("%w: cannot delete device in use", ErrInvalidDelete)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete device: %w", err)
	}

	return nil
}

// internal helper
func validateState(s string) error {
	switch model.DeviceState(s) {
	case model.StateAvailable, model.StateInUse, model.StateInactive:
		return nil
	default:
		return ErrInvalidState
	}
}
