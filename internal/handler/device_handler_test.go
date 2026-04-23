package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dev.paultergust/devices-api/internal/dto"
	"dev.paultergust/devices-api/internal/model"
	"dev.paultergust/devices-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockDeviceService struct {
	createFn func(ctx context.Context, req dto.CreateDeviceRequest) (*model.Device, error)
	getFn    func(ctx context.Context, id uuid.UUID) (*model.Device, error)
	listFn   func(ctx context.Context, brand, state *string) ([]model.Device, error)
	updateFn func(ctx context.Context, id uuid.UUID, req dto.UpdateDeviceRequest) (*model.Device, error)
	deleteFn func(ctx context.Context, id uuid.UUID) error
}

func (m *mockDeviceService) Create(ctx context.Context, req dto.CreateDeviceRequest) (*model.Device, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return nil, nil
}

func (m *mockDeviceService) Get(ctx context.Context, id uuid.UUID) (*model.Device, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return nil, nil
}

func (m *mockDeviceService) List(ctx context.Context, brand, state *string) ([]model.Device, error) {
	if m.listFn != nil {
		return m.listFn(ctx, brand, state)
	}
	return nil, nil
}

func (m *mockDeviceService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateDeviceRequest) (*model.Device, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, req)
	}
	return nil, nil
}

func (m *mockDeviceService) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.POST("/devices", h.CreateDevice)
	r.GET("/devices", h.ListDevices)
	r.GET("/devices/:id", h.GetDevice)
	r.PATCH("/devices/:id", h.UpdateDevice)
	r.DELETE("/devices/:id", h.DeleteDevice)
	return r
}

func TestCreateDevice_Success(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	svc := &mockDeviceService{
		createFn: func(ctx context.Context, req dto.CreateDeviceRequest) (*model.Device, error) {
			if req.Name != "iPhone" {
				t.Fatalf("expected name iPhone, got %s", req.Name)
			}
			return &model.Device{
				ID:        id,
				Name:      req.Name,
				Brand:     req.Brand,
				State:     model.StateAvailable,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	h := NewHandler(svc)
	r := setupRouter(h)

	body := `{"name":"iPhone","brand":"Apple","state":"available"}`
	req := httptest.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusCreated, w.Code, w.Body.String())
	}

	var resp dto.DeviceResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.ID != id.String() {
		t.Fatalf("expected id %s, got %s", id.String(), resp.ID)
	}
	if resp.Name != "iPhone" {
		t.Fatalf("expected name iPhone, got %s", resp.Name)
	}
}

func TestCreateDevice_InvalidJSON(t *testing.T) {
	t.Parallel()

	h := NewHandler(&mockDeviceService{})
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(`{"name":`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateDevice_InvalidState(t *testing.T) {
	t.Parallel()

	h := NewHandler(&mockDeviceService{})
	r := setupRouter(h)

	body := `{"name":"iPhone","brand":"Apple","state":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestGetDevice_Success(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	svc := &mockDeviceService{
		getFn: func(ctx context.Context, gotID uuid.UUID) (*model.Device, error) {
			if gotID != id {
				t.Fatalf("expected id %s, got %s", id, gotID)
			}
			return &model.Device{
				ID:        id,
				Name:      "Pixel",
				Brand:     "Google",
				State:     model.StateAvailable,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	h := NewHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/devices/"+id.String(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetDevice_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewHandler(&mockDeviceService{})
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/devices/not-a-uuid", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetDevice_NotFound(t *testing.T) {
	t.Parallel()

	svc := &mockDeviceService{
		getFn: func(ctx context.Context, id uuid.UUID) (*model.Device, error) {
			return nil, service.ErrNot[118;1:3uFound
		},
	}

	h := NewHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/devices/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestListDevices_Success(t *testing.T) {
	t.Parallel()

	svc := &mockDeviceService{
		listFn: func(ctx context.Context, brand, state *string) ([]model.Device, error) {
			if brand == nil || *brand != "Apple" {
				t.Fatalf("expected brand Apple, got %v", brand)
			}
			if state == nil || *state != "available" {
				t.Fatalf("expected state available, got %v", state)
			}
			return []model.Device{
				{
					ID:        uuid.New(),
					Name:      "iPhone",
					Brand:     "Apple",
					State:     model.StateAvailable,
					CreatedAt: time.Now(),
				},
			}, nil
		},
	}

	h := NewHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/devices?brand=Apple&state=available", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp []dto.DeviceResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 device, got %d", len(resp))
	}
}

func TestListDevices_InternalError(t *testing.T) {
	t.Parallel()

	svc := &mockDeviceService{
		listFn: func(ctx context.Context, brand, state *string) ([]model.Device, error) {
			return nil, errors.New("db exploded")
		},
	}

	h := NewHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUpdateDevice_Success(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	svc := &mockDeviceService{
		updateFn: func(ctx context.Context, gotID uuid.UUID, req dto.UpdateDeviceRequest) (*model.Device, error) {
			if gotID != id {
				t.Fatalf("expected id %s, got %s", id, gotID)
			}
			if req.Name == nil || *req.Name != "Updated" {
				t.Fatalf("expected updated name, got %+v", req.Name)
			}
			return &model.Device{
				ID:        id,
				Name:      "Updated",
				Brand:     "Apple",
				State:     model.StateAvailable,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	h := NewHandler(svc)
	r := setupRouter(h)

	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPatch, "/devices/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestUpdateDevice_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewHandler(&mockDeviceService{})
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPatch, "/devices/not-a-uuid", bytes.NewBufferString(`{"name":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateDevice_InvalidJSON(t *testing.T) {
	t.Parallel()

	h := NewHandler(&mockDeviceService{})
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPatch, "/devices/"+uuid.New().String(), bytes.NewBufferString(`{"name":`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateDevice_InvalidState(t *testing.T) {
	t.Parallel()

	h := NewHandler(&mockDeviceService{})
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPatch, "/devices/"+uuid.New().String(), bytes.NewBufferString(`{"state":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestUpdateDevice_ServiceInvalidUpdate(t *testing.T) {
	t.Parallel()

	svc := &mockDeviceService{
		updateFn: func(ctx context.Context, id uuid.UUID, req dto.UpdateDeviceRequest) (*model.Device, error) {
			return nil, service.ErrInvalidUpdate
		},
	}

	h := NewHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPatch, "/devices/"+uuid.New().String(), bytes.NewBufferString(`{"name":"new"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestDeleteDevice_Success(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	svc := &mockDeviceService{
		deleteFn: func(ctx context.Context, gotID uuid.UUID) error {
			if gotID != id {
				t.Fatalf("expected id %s, got %s", id, gotID)
			}
			return nil
		},
	}

	h := NewHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/devices/"+id.String(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestDeleteDevice_InvalidID(t *testing.T) {
	t.Parallel()

	h := NewHandler(&mockDeviceService{})
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/devices/not-a-uuid", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteDevice_NotFound(t *testing.T) {
	t.Parallel()

	svc := &mockDeviceService{
		deleteFn: func(ctx context.Context, id uuid.UUID) error {
			return service.ErrNotFound
		},
	}

	h := NewHandler(svc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/devices/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusNotFound, w.Code, w.Body.String())
	}
}
