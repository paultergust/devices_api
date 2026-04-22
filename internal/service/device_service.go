package service

import (
    "context"
    "errors"
    "time"

    "dev.paultergust/devices-api/internal/dto"
    "dev.paultergust/devices-api/internal/model"
    "dev.paultergust/devices-api/internal/repository"

    "github.com/google/uuid"
)

type DeviceService struct {
    repo *repository.DeviceRepository
}

func NewDeviceService(r *repository.DeviceRepository) *DeviceService {
    return &DeviceService{repo: r}
}

func (s *DeviceService) Create(ctx context.Context, req dto.CreateDeviceRequest) (*model.Device, error) {
    d := &model.Device{
        ID:        uuid.New(),
        Name:      req.Name,
        Brand:     req.Brand,
        State:     model.DeviceState(req.State),
        CreatedAt: time.Now(),
    }

    if err := s.repo.Create(ctx, d); err != nil {
        return nil, err
    }

    return d, nil
}

func (s *DeviceService) Get(ctx context.Context, id uuid.UUID) (*model.Device, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *DeviceService) List(ctx context.Context, brand, state *string) ([]model.Device, error) {
    return s.repo.List(ctx, brand, state)
}

func (s *DeviceService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateDeviceRequest) (*model.Device, error) {
    d, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    if d.State == model.StateInUse {
        if req.Name != nil || req.Brand != nil {
            return nil, errors.New("cannot update name or brand when device is in use")
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
        return nil, err
    }

    return d, nil
}

func (s *DeviceService) Delete(ctx context.Context, id uuid.UUID) error {
    d, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    if d.State == model.StateInUse {
        return errors.New("cannot delete device in use")
    }

    return s.repo.Delete(ctx, id)
}
