package dto

import (
	"errors"
	"time"

	"dev.paultergust/devices-api/internal/model"
)

var (
	ErrInvalidState = errors.New("invalid device state")
	ErrNilDevice    = errors.New("device is nil")
)

type CreateDeviceRequest struct {
	Name  string `json:"name" binding:"required"`
	Brand string `json:"brand" binding:"required"`
	State string `json:"state" binding:"required,oneof=available in-use inactive"`
}

type UpdateDeviceRequest struct {
	Name  *string `json:"name,omitempty"`
	Brand *string `json:"brand,omitempty"`
	State *string `json:"state,omitempty" binding:"omitempty,oneof=available in-use inactive"`
}

type DeviceResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Brand     string    `json:"brand"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

func (r CreateDeviceRequest) Validate() error {
	if !isValidState(r.State) {
		return ErrInvalidState
	}
	return nil
}

func (r UpdateDeviceRequest) Validate() error {
	if r.State != nil && !isValidState(*r.State) {
		return ErrInvalidState
	}
	return nil
}

func isValidState(s string) bool {
	switch model.DeviceState(s) {
	case model.StateAvailable, model.StateInUse, model.StateInactive:
		return true
	default:
		return false
	}
}

func ToDeviceResponse(d *model.Device) (DeviceResponse, error) {
	if d == nil {
		return DeviceResponse{}, ErrNilDevice
	}

	return DeviceResponse{
		ID:        d.ID.String(),
		Name:      d.Name,
		Brand:     d.Brand,
		State:     string(d.State),
		CreatedAt: d.CreatedAt,
	}, nil
}

func ToDeviceResponses(devices []model.Device) ([]DeviceResponse, error) {
	res := make([]DeviceResponse, 0, len(devices))

	for i := range devices {
		r, err := ToDeviceResponse(&devices[i])
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return res, nil
}
