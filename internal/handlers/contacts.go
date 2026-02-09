package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/services"
)

type ContactHandler struct {
	contactService *services.ContactService
}

func (h *ContactHandler) GetAllByClientID(w http.ResponseWriter, r *http.Request) {
	clientID := chi.URLParam(r, "clientID")

	if clientID == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(clientID)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	allContacts, err := h.contactService.GetAllByClientID(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, allContacts)
}

func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var body models.Contact

	if !decodeJSONBody(w, r, &body) {
		return
	}

	err := h.contactService.CreateContact(r.Context(), body)
	if err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")

	if ID == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(ID)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	err = h.contactService.DeleteContact(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, ID)
}

func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")

	if ID == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(ID)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	var updates models.ContactUpdate

	if !decodeJSONBody(w, r, &updates) {
		return
	}

	updated, err := h.contactService.UpdateContact(r.Context(), uuid, updates)
	if err != nil {
		serverError(w, err)
		return
	}

	if updated == nil {
		notFound(w)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}
