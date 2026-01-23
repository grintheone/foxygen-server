package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/services"
)

type AgreementHandler struct {
	service *services.AgreementService
}

func (h *AgreementHandler) GetAgreementsByField(w http.ResponseWriter, r *http.Request) {
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

	field := r.URL.Query().Get("field")
	if field == "" {
		clientError(w, http.StatusBadRequest)
		return
	}
	active := r.URL.Query().Get("active")
	if active == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	agreements, err := h.service.GetAgreementsByField(r.Context(), uuid, field, active)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, agreements)
}
