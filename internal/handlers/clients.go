package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/services"
)

type ClientHandler struct {
	clientService *services.ClientService
}

func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.clientService.ListClients(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, clients)
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var payload models.Client

	if !decodeJSONBody(w, r, &payload) {
		return
	}

	created, err := h.clientService.CreateClient(r.Context(), payload)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
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

	var payload models.ClientUpdate

	if !decodeJSONBody(w, r, &payload) {
		return
	}

	updated, err := h.clientService.UpdateClient(r.Context(), uuid, payload)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
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

	err = h.clientService.DeleteClient(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, uuid)
}

func (h *ClientHandler) GetClientByID(w http.ResponseWriter, r *http.Request) {
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

	client, err := h.clientService.GetClientByID(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, client)
}
