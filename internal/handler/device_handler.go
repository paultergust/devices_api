package handler

import (
	"context"
	"errors"
	"net/http"

	"dev.paultergust/devices-api/internal/dto"
	"dev.paultergust/devices-api/internal/model"
	"dev.paultergust/devices-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)
type deviceService interface {
	Create(ctx context.Context, req dto.CreateDeviceRequest) (*model.Device, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Device, error)
	List(ctx context.Context, brand, state *string) ([]model.Device, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateDeviceRequest) (*model.Device, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	svc deviceService
}

func NewHandler(s deviceService) *Handler {
	return &Handler{svc: s}
}


// central error → HTTP mapping
func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, dto.ErrInvalidState):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
	case errors.Is(err, service.ErrInvalidUpdate):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		// do not leak internal errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func (h *Handler) CreateDevice(c *gin.Context) {
	var req dto.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		handleError(c, err)
		return
	}

	d, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	resp, err := dto.ToDeviceResponse(d)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	d, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	resp, err := dto.ToDeviceResponse(d)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListDevices(c *gin.Context) {
	brand := c.Query("brand")
	state := c.Query("state")

	var brandPtr, statePtr *string

	if brand != "" {
		brandPtr = &brand
	}
	if state != "" {
		statePtr = &state
	}

	devices, err := h.svc.List(c.Request.Context(), brandPtr, statePtr)
	if err != nil {
		handleError(c, err)
		return
	}

	resp, err := dto.ToDeviceResponses(devices)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		handleError(c, err)
		return
	}

	d, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		handleError(c, err)
		return
	}

	resp, err := dto.ToDeviceResponse(d)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
