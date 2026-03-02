package handlers

import (
	"net/http"

	"github.com/grintheone/foxygen-server/internal/services"
)

type ManufacturerHandler struct {
	service *services.ManufacturerService
}

func (h *ManufacturerHandler) ListAllManufacturers(w http.ResponseWriter, r *http.Request) {
	manufacturers, err := h.service.ListAllManufacturers(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, manufacturers)
}
