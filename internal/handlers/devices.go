package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/services"
)

type DeviceHandler struct {
	deviceService *services.DeviceService
}

func (h *DeviceHandler) GetAllDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.deviceService.GetAllDevices(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, devices)
}

func (h *DeviceHandler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")

	if uuidStr == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	device, err := h.deviceService.GetDeviceByID(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, device)
}

func (h *DeviceHandler) RemoveDeviceByID(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")

	if uuidStr == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	err = h.deviceService.RemoveDeviceByID(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, uuid)
}

func (h *DeviceHandler) CreateNewDevice(w http.ResponseWriter, r *http.Request) {
	var body models.Device

	if !decodeJSONBody(w, r, &body) {
		return
	}

	device, err := h.deviceService.CreateNewDevice(r.Context(), body)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, device)
}

func (h *DeviceHandler) UpdateDeviceByID(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")

	if uuidStr == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	var body models.DeviceUpdates

	if !decodeJSONBody(w, r, &body) {
		return
	}

	updated, err := h.deviceService.UpdateDeviceByID(r.Context(), uuid, body)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}
