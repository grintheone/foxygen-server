package handlers

import (
	"net/http"

	"github.com/grintheone/foxygen-server/internal/services"
)

type ResearchTypeHandler struct {
	service *services.ResearchTypeService
}

func (h *ResearchTypeHandler) ListAllResearchTypes(w http.ResponseWriter, r *http.Request) {
	researchTypes, err := h.service.ListAllResearchTypes(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, researchTypes)
}
