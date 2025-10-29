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

func (h *AgreementHandler) GetClientAgreements(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "clientID")

	if uuidStr == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	agreements, err := h.service.GetClientAgreements(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, agreements)
}

func (h *AgreementHandler) GetAgreementsByField(w http.ResponseWriter, r *http.Request) {
	field := chi.URLParam(r, "field")
	if field == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

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

	agreements, err := h.service.GetAgreementsByField(r.Context(), field, uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, agreements)
}
