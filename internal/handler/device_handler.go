package handler

import (
    "net/http"

    "dev.paultergust/devices-api/internal/dto"
    "dev.paultergust/devices-api/internal/service"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type Handler struct {
    svc *service.DeviceService
}

func NewHandler(s *service.DeviceService) *Handler {
    return &Handler{svc: s}
}

func (h *Handler) CreateDevice(c *gin.Context) {
    var req dto.CreateDeviceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    d, err := h.svc.Create(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, d)
}

func (h *Handler) GetDevice(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    d, err := h.svc.Get(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    c.JSON(http.StatusOK, d)
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, devices)
}

func (h *Handler) UpdateDevice(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var req dto.UpdateDeviceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    d, err := h.svc.Update(c.Request.Context(), id, req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, d)
}

func (h *Handler) DeleteDevice(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    err = h.svc.Delete(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}
