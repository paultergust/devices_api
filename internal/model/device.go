package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type DeviceState string

const (
	StateAvailable DeviceState = "available"
	StateInUse     DeviceState = "in-use"
	StateInactive  DeviceState = "inactive"
)

var (
	ErrInvalidDeviceState = errors.New("invalid device state")
)

func (s DeviceState) IsValid() bool {
	switch s {
	case StateAvailable, StateInUse, StateInactive:
		return true
	default:
		return false
	}
}

func (s DeviceState) Validate() error {
	if !s.IsValid() {
		return ErrInvalidDeviceState
	}
	return nil
}

type Device struct {
	ID        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Brand     string      `json:"brand"`
	State     DeviceState `json:"state"`
	CreatedAt time.Time   `json:"created_at"`
}

// Constructor enforces invariants
func NewDevice(name, brand string, state DeviceState) (*Device, error) {
	if err := state.Validate(); err != nil {
		return nil, err
	}

	return &Device{
		ID:        uuid.New(),
		Name:      name,
		Brand:     brand,
		State:     state,
		CreatedAt: time.Now(),
	}, nil
}

// Domain helpers (avoid string comparisons everywhere)
func (d *Device) IsInUse() bool {
	return d.State == StateInUse
}

func (d *Device) CanUpdateDetails() bool {
	return d.State != StateInUse
}

func (d *Device) CanDelete() bool {
	return d.State != StateInUse
}
