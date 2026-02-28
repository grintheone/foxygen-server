package handlers

import (
	"net/http"

	"github.com/grintheone/foxygen-server/internal/services"
)

type RegionHandler struct {
	service *services.RegionService
}

func (h *RegionHandler) ListAllRegions(w http.ResponseWriter, r *http.Request) {
	regions, err := h.service.ListAllRegions(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, regions)
}
