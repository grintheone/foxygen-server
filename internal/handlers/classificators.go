package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/services"
)

type ClassificatorHandler struct {
	classificatorService *services.ClassificatorService
}

func (h *ClassificatorHandler) GetClassificatorByID(w http.ResponseWriter, r *http.Request) {
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

	device, err := h.classificatorService.GetClassificatorByID(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, device)
}

func (h *ClassificatorHandler) GetDevicesByClassificatorID(w http.ResponseWriter, r *http.Request) {
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

	devices, err := h.classificatorService.GetDevicesByClassificatorID(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, devices)
}

func (h *ClassificatorHandler) NewClassificator(w http.ResponseWriter, r *http.Request) {
	var body models.Classificator

	if !decodeJSONBody(w, r, &body) {
		return
	}

	err := h.classificatorService.NewClassificator(r.Context(), body)
	if err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ClassificatorHandler) RemoveClassificatorByID(w http.ResponseWriter, r *http.Request) {
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

	err = h.classificatorService.RemoveClassificatorByID(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, uuid)
}

func (h *ClassificatorHandler) UpdateClassificatorInfo(w http.ResponseWriter, r *http.Request) {
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

	var body models.ClassificatorUpdate

	if !decodeJSONBody(w, r, &body) {
		return
	}

	updated, err := h.classificatorService.UpdateClassificatorInfo(r.Context(), uuid, body)

	writeJSON(w, http.StatusOK, updated)
}
