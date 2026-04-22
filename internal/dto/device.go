package dto

import (
	"time"

	"dev.paultergust/devices-api/internal/model"
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

func ToDeviceResponse(d *model.Device) DeviceResponse {
    return DeviceResponse{
        ID:        d.ID.String(),
        Name:      d.Name,
        Brand:     d.Brand,
        State:     string(d.State),
        CreatedAt: d.CreatedAt,
    }
}
