package repository

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"dev.paultergust/devices-api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestRepository(t *testing.T) (*DeviceRepository, *pgxpool.Pool) {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("failed to ping database: %v", err)
	}

	if _, err := pool.Exec(ctx, `TRUNCATE TABLE devices RESTART IDENTITY`); err != nil {
		pool.Close()
		t.Fatalf("failed to truncate devices table: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return NewDeviceRepository(pool), pool
}

func seedDevice(t *testing.T, repo *DeviceRepository, name, brand string, state model.DeviceState) *model.Device {
	t.Helper()

	device := &model.Device{
		ID:        uuid.New(),
		Name:      name,
		Brand:     brand,
		State:     state,
		CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
	}

	if err := repo.Create(context.Background(), device); err != nil {
		t.Fatalf("failed to seed device: %v", err)
	}

	return device
}

func TestDeviceRepository_CreateAndGetByID(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	device := &model.Device{
		ID:        uuid.New(),
		Name:      "iPhone 15",
		Brand:     "Apple",
		State:     model.StateAvailable,
		CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
	}

	if err := repo.Create(ctx, device); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := repo.GetByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if got.ID != device.ID {
		t.Fatalf("expected ID %v, got %v", device.ID, got.ID)
	}
	if got.Name != device.Name {
		t.Fatalf("expected Name %q, got %q", device.Name, got.Name)
	}
	if got.Brand != device.Brand {
		t.Fatalf("expected Brand %q, got %q", device.Brand, got.Brand)
	}
	if got.State != device.State {
		t.Fatalf("expected State %q, got %q", device.State, got.State)
	}
}

func TestDeviceRepository_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeviceRepository_List_All(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepository(t)

	seedDevice(t, repo, "iPhone", "Apple", model.StateAvailable)
	seedDevice(t, repo, "Galaxy", "Samsung", model.StateInactive)

	got, err := repo.List(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(got))
	}
}

func TestDeviceRepository_List_ByBrand(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepository(t)

	seedDevice(t, repo, "iPhone", "Apple", model.StateAvailable)
	seedDevice(t, repo, "Galaxy", "Samsung", model.StateAvailable)

	brand := "Apple"
	got, err := repo.List(context.Background(), &brand, nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 device, got %d", len(got))
	}
	if got[0].Brand != "Apple" {
		t.Fatalf("expected brand Apple, got %q", got[0].Brand)
	}
}

func TestDeviceRepository_List_ByState(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepository(t)

	seedDevice(t, repo, "iPhone", "Apple", model.StateAvailable)
	seedDevice(t, repo, "Galaxy", "Samsung", model.StateInactive)

	state := "inactive"
	got, err := repo.List(context.Background(), nil, &state)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 device, got %d", len(got))
	}
	if got[0].State != model.StateInactive {
		t.Fatalf("expected state inactive, got %q", got[0].State)
	}
}

func TestDeviceRepository_List_ByBrandAndState(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepository(t)

	seedDevice(t, repo, "iPhone", "Apple", model.StateAvailable)
	seedDevice(t, repo, "Old iPhone", "Apple", model.StateInactive)
	seedDevice(t, repo, "Galaxy", "Samsung", model.StateAvailable)

	brand := "Apple"
	state := "available"

	got, err := repo.List(context.Background(), &brand, &state)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 device, got %d", len(got))
	}
	if got[0].Brand != "Apple" || got[0].State != model.StateAvailable {
		t.Fatalf("unexpected device returned: %+v", got[0])
	}
}

func TestDeviceRepository_Update_Success(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	device := seedDevice(t, repo, "iPhone", "Apple", model.StateAvailable)

	device.Name = "iPhone Pro"
	device.Brand = "Apple Updated"
	device.State = model.StateInactive

	if err := repo.Update(ctx, device); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	got, err := repo.GetByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if got.Name != "iPhone Pro" {
		t.Fatalf("expected updated name, got %q", got.Name)
	}
	if got.Brand != "Apple Updated" {
		t.Fatalf("expected updated brand, got %q", got.Brand)
	}
	if got.State != model.StateInactive {
		t.Fatalf("expected updated state inactive, got %q", got.State)
	}
}

func TestDeviceRepository_Update_NotFound(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepository(t)

	device := &model.Device{
		ID:        uuid.New(),
		Name:      "Missing",
		Brand:     "NoBrand",
		State:     model.StateAvailable,
		CreatedAt: time.Now().UTC(),
	}

	err := repo.Update(context.Background(), device)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeviceRepository_Delete_Success(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepository(t)
	ctx := context.Background()

	device := seedDevice(t, repo, "iPhone", "Apple", model.StateAvailable)

	if err := repo.Delete(ctx, device.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err := repo.GetByID(ctx, device.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeviceRepository_Delete_NotFound(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepository(t)

	err := repo.Delete(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
