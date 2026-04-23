package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"dev.paultergust/devices-api/internal/dto"
	"dev.paultergust/devices-api/internal/model"
	"dev.paultergust/devices-api/internal/repository"

	"github.com/google/uuid"
)

type mockDeviceRepository struct {
	createFn  func(ctx context.Context, d *model.Device) error
	getByIDFn func(ctx context.Context, id uuid.UUID) (*model.Device, error)
	listFn    func(ctx context.Context, brand, state *string) ([]model.Device, error)
	updateFn  func(ctx context.Context, d *model.Device) error
	deleteFn  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockDeviceRepository) Create(ctx context.Context, d *model.Device) error {
	if m.createFn != nil {
		return m.createFn(ctx, d)
	}
	return nil
}

func (m *mockDeviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Device, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockDeviceRepository) List(ctx context.Context, brand, state *string) ([]model.Device, error) {
	if m.listFn != nil {
		return m.listFn(ctx, brand, state)
	}
	return nil, nil
}

func (m *mockDeviceRepository) Update(ctx context.Context, d *model.Device) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, d)
	}
	return nil
}

func (m *mockDeviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func TestDeviceService_Create_Success(t *testing.T) {
	t.Parallel()

	req := dto.CreateDeviceRequest{
		Name:  "iPhone 15",
		Brand: "Apple",
		State: "available",
	}

	repo := &mockDeviceRepository{
		createFn: func(ctx context.Context, d *model.Device) error {
			if d.ID == uuid.Nil {
				t.Fatal("expected generated UUID")
			}
			if d.Name != req.Name {
				t.Fatalf("expected name %q, got %q", req.Name, d.Name)
			}
			if d.Brand != req.Brand {
				t.Fatalf("expected brand %q, got %q", req.Brand, d.Brand)
			}
			if d.State != model.StateAvailable {
				t.Fatalf("expected state %q, got %q", model.StateAvailable, d.State)
			}
			if d.CreatedAt.IsZero() {
				t.Fatal("expected CreatedAt to be set")
			}
			return nil
		},
	}

	svc := NewDeviceService(repo)

	got, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got == nil {
		t.Fatal("expected device, got nil")
	}
	if got.Name != req.Name || got.Brand != req.Brand || got.State != model.StateAvailable {
		t.Fatalf("unexpected device: %+v", got)
	}
}

func TestDeviceService_Create_InvalidState(t *testing.T) {
	t.Parallel()

	repo := &mockDeviceRepository{}
	svc := NewDeviceService(repo)

	_, err := svc.Create(context.Background(), dto.CreateDeviceRequest{
		Name:  "Device",
		Brand: "Brand",
		State: "broken",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState, got %v", err)
	}
}

func TestDeviceService_Create_RepositoryError(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("db down")
	repo := &mockDeviceRepository{
		createFn: func(ctx context.Context, d *model.Device) error {
			return repoErr
		},
	}

	svc := NewDeviceService(repo)

	_, err := svc.Create(context.Background(), dto.CreateDeviceRequest{
		Name:  "Galaxy",
		Brand: "Samsung",
		State: "available",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
}

func TestDeviceService_Get_Success(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	expected := &model.Device{
		ID:        id,
		Name:      "Pixel",
		Brand:     "Google",
		State:     model.StateAvailable,
		CreatedAt: time.Now(),
	}

	repo := &mockDeviceRepository{
		getByIDFn: func(ctx context.Context, gotID uuid.UUID) (*model.Device, error) {
			if gotID != id {
				t.Fatalf("expected id %v, got %v", id, gotID)
			}
			return expected, nil
		},
	}

	svc := NewDeviceService(repo)

	got, err := svc.Get(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != expected {
		t.Fatalf("expected %+v, got %+v", expected, got)
	}
}

func TestDeviceService_Get_NotFound(t *testing.T) {
	t.Parallel()

	repo := &mockDeviceRepository{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Device, error) {
			return nil, repository.ErrNotFound
		},
	}

	svc := NewDeviceService(repo)

	_, err := svc.Get(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeviceService_List_Success(t *testing.T) {
	t.Parallel()

	brand := "Apple"
	state := "available"

	expected := []model.Device{
		{
			ID:        uuid.New(),
			Name:      "iPhone",
			Brand:     "Apple",
			State:     model.StateAvailable,
			CreatedAt: time.Now(),
		},
	}

	repo := &mockDeviceRepository{
		listFn: func(ctx context.Context, gotBrand, gotState *string) ([]model.Device, error) {
			if gotBrand == nil || *gotBrand != brand {
				t.Fatalf("expected brand filter %q, got %v", brand, gotBrand)
			}
			if gotState == nil || *gotState != state {
				t.Fatalf("expected state filter %q, got %v", state, gotState)
			}
			return expected, nil
		},
	}

	svc := NewDeviceService(repo)

	got, err := svc.List(context.Background(), &brand, &state)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != 1 || got[0].Name != "iPhone" {
		t.Fatalf("unexpected devices: %+v", got)
	}
}

func TestDeviceService_List_InvalidState(t *testing.T) {
	t.Parallel()

	state := "bad-state"
	repo := &mockDeviceRepository{}
	svc := NewDeviceService(repo)

	_, err := svc.List(context.Background(), nil, &state)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected ErrInvalidState, got %v", err)
	}
}

func TestDeviceService_Update_Success(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	newName := "Updated Device"
	newBrand := "Updated Brand"
	newState := "inactive"

	existing := &model.Device{
		ID:        id,
		Name:      "Old Device",
		Brand:     "Old Brand",
		State:     model.StateAvailable,
		CreatedAt: time.Now(),
	}

	repo := &mockDeviceRepository{
		getByIDFn: func(ctx context.Context, gotID uuid.UUID) (*model.Device, error) {
			return existing, nil
		},
		updateFn: func(ctx context.Context, d *model.Device) error {
			if d.Name != newName {
				t.Fatalf("expected name %q, got %q", newName, d.Name)
			}
			if d.Brand != newBrand {
				t.Fatalf("expected brand %q, got %q", newBrand, d.Brand)
			}
			if d.State != model.StateInactive {
				t.Fatalf("expected state %q, got %q", model.StateInactive, d.State)
			}
			return nil
		},
	}

	svc := NewDeviceService(repo)

	got, err := svc.Update(context.Background(), id, dto.UpdateDeviceRequest{
		Name:  &newName,
		Brand: &newBrand,
		State: &newState,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Name != newName || got.Brand != newBrand || got.State != model.StateInactive {
		t.Fatalf("unexpected updated device: %+v", got)
	}
}

func TestDeviceService_Update_InUseCannotChangeNameOrBrand(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	newName := "New Name"

	repo := &mockDeviceRepository{
		getByIDFn: func(ctx context.Context, gotID uuid.UUID) (*model.Device, error) {
			return &model.Device{
				ID:        id,
				Name:      "Current",
				Brand:     "Brand",
				State:     model.StateInUse,
				CreatedAt: time.Now(),
			}, nil
		},
		updateFn: func(ctx context.Context, d *model.Device) error {
			t.Fatal("update should not be called")
[118;1:3u			return nil
		},
	}

	svc := NewDeviceService(repo)

	_, err := svc.Update(context.Background(), id, dto.UpdateDeviceRequest{
		Name: &newName,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidUpdate) {
		t.Fatalf("expected ErrInvalidUpdate, got %v", err)
	}
}

func TestDeviceService_Update_InUseCanChangeState(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	newState := "inactive"

	repo := &mockDeviceRepository{
		getByIDFn: func(ctx context.Context, gotID uuid.UUID) (*model.Device, error) {
			return &model.Device{
				ID:        id,
				Name:      "Current",
				Brand:     "Brand",
				State:     model.StateInUse,
				CreatedAt: time.Now(),
			}, nil
		},
		updateFn: func(ctx context.Context, d *model.Device) error {
			if d.State != model.StateInactive {
				t.Fatalf("expected state to change to inactive, got %q", d.State)
			}
			return nil
		},
	}

	svc := NewDeviceService(repo)

	got, err := svc.Update(context.Background(), id, dto.UpdateDeviceRequest{
		State: &newState,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.State != model.StateInactive {
		t.Fatalf("expected updated state inactive, got %q", got.State)
	}
}

func TestDeviceService_Update_NotFound(t *testing.T) {
	t.Parallel()

	repo := &mockDeviceRepository{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Device, error) {
			return nil, repository.ErrNotFound
		},
	}

	svc := NewDeviceService(repo)

	_, err := svc.Update(context.Background(), uuid.New(), dto.UpdateDeviceRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeviceService_Delete_Success(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	repo := &mockDeviceRepository{
		getByIDFn: func(ctx context.Context, gotID uuid.UUID) (*model.Device, error) {
			return &model.Device{
				ID:        id,
				Name:      "Device",
				Brand:     "Brand",
				State:     model.StateAvailable,
				CreatedAt: time.Now(),
			}, nil
		},
		deleteFn: func(ctx context.Context, gotID uuid.UUID) error {
			if gotID != id {
				t.Fatalf("expected id %v, got %v", id, gotID)
			}
			return nil
		},
	}

	svc := NewDeviceService(repo)

	if err := svc.Delete(context.Background(), id); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeviceService_Delete_InUse(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	repo := &mockDeviceRepository{
		getByIDFn: func(ctx context.Context, gotID uuid.UUID) (*model.Device, error) {
			return &model.Device{
				ID:        id,
				Name:      "Device",
				Brand:     "Brand",
				State:     model.StateInUse,
				CreatedAt: time.Now(),
			}, nil
		},
		deleteFn: func(ctx context.Context, gotID uuid.UUID) error {
			t.Fatal("delete should not be called")
			return nil
		},
	}

	svc := NewDeviceService(repo)

	err := svc.Delete(context.Background(), id)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidDelete) {
		t.Fatalf("expected ErrInvalidDelete, got %v", err)
	}
}

func TestDeviceService_Delete_NotFound(t *testing.T) {
	t.Parallel()

	repo := &mockDeviceRepository{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*model.Device, error) {
			return nil, repository.ErrNotFound
		},
	}

	svc := NewDeviceService(repo)

	err := svc.Delete(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
