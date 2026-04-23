package dto

import (
	"errors"
	"testing"
	"time"

	"dev.paultergust/devices-api/internal/model"

	"github.com/google/uuid"
)

func TestCreateDeviceRequest_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     CreateDeviceRequest
		wantErr bool
	}{
		{
			name: "valid state",
			req: CreateDeviceRequest{
				Name:  "iPhone",
				Brand: "Apple",
				State: "available",
			},
			wantErr: false,
		},
		{
			name: "invalid state",
			req: CreateDeviceRequest{
				Name:  "iPhone",
				Brand: "Apple",
				State: "broken",
			},
			wantErr: true,
		},
		{
			name: "empty state",
			req: CreateDeviceRequest{
				Name:  "iPhone",
				Brand: "Apple",
				State: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.req.Validate()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, ErrInvalidState) {
					t.Fatalf("expected ErrInvalidState, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestUpdateDeviceRequest_Validate(t *testing.T) {
	t.Parallel()

	valid := "available"
	invalid := "broken"

	tests := []struct {
		name    string
		req     UpdateDeviceRequest
		wantErr bool
	}{
		{
			name:    "nil state",
			req:     UpdateDeviceRequest{},
			wantErr: false,
		},
		{
			name: "valid state",
			req: UpdateDeviceRequest{
				State: &valid,
			},
			wantErr: false,
		},
		{
			name: "invalid state",
			req: UpdateDeviceRequest{
				State: &invalid,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.req.Validate()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, ErrInvalidState) {
					t.Fatalf("expected ErrInvalidState, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestToDeviceResponse(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	now := time.Now()

	device := &model.Device{
		ID:        id,
		Name:      "Pixel",
		Brand:     "Google",
		State:     model.StateAvailable,
		CreatedAt: now,
	}

	resp, err := ToDeviceResponse(device)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.ID != id.String() {
		t.Fatalf("expected id %s, got %s", id.String(), resp.ID)
	}
	if resp.Name != "Pixel" {
		t.Fatalf("expected name Pixel, got %s", resp.Name)
	}
	if resp.Brand != "Google" {
		t.Fatalf("expected brand Google, got %s", resp.Brand)
	}
	if resp.State != "available" {
		t.Fatalf("expected state available, got %s", resp.State)
	}
	if !resp.CreatedAt.Equal(now) {
		t.Fatalf("expected CreatedAt %v, got %v", now, resp.CreatedAt)
	}
}

func TestToDeviceResponse_Nil(t *testing.T) {
	t.Parallel()

	_, err := ToDeviceResponse(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNilDevice) {
		t.Fatalf("expected ErrNilDevice, got %v", err)
	}
}

func TestToDeviceResponses(t *testing.T) {
	t.Parallel()

	devices := []model.Device{
		{
			ID:        uuid.New(),
			Name:      "iPhone",
			Brand:     "Apple",
			State:     model.StateAvailable,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Galaxy",
			Brand:     "Samsung",
			State:     model.StateInactive,
			CreatedAt: time.Now(),
		},
	}

	resp, err := ToDeviceResponses(devices)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp) != len(devices) {
		t.Fatalf("expected %d devices, got %d", len(devices), len(resp))
	}

	for i := range resp {
		if resp[i].ID != devices[i].ID.String() {
			t.Fatalf("mismatch at index %d", i)
		}
	}
}

func TestToDeviceResponses_Empty(t *testing.T) {
	t.Parallel()

	resp, err := ToDeviceResponses([]model.Device{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp) != 0 {
		t.Fatalf("expected empty slice, got %d", len(resp))
	}
}
