package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/middlewares"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/services"
)

type TicketHandler struct {
	ticketService *services.TicketService
}

func (h *TicketHandler) ListAllTickets(w http.ResponseWriter, r *http.Request) {
	executorID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		serverError(w, fmt.Errorf("Unable to check for user role"))
		return
	}

	tickets, err := h.ticketService.ListAllTickets(r.Context(), executorID)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tickets)
}

func (h *TicketHandler) GetTicketByID(w http.ResponseWriter, r *http.Request) {
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

	ticket, err := h.ticketService.GetTicketByID(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) DeleteTicketByID(w http.ResponseWriter, r *http.Request) {
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

	err = h.ticketService.DeleteTicketByID(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, uuid)
}

func (h *TicketHandler) CreateNewTicket(w http.ResponseWriter, r *http.Request) {
	var body models.RawTicket

	if !decodeJSONBody(w, r, &body) {
		return
	}

	created, err := h.ticketService.CreateNewTicket(r.Context(), body)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *TicketHandler) UpdateTicketInfo(w http.ResponseWriter, r *http.Request) {
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

	var updates models.TicketUpdates

	if !decodeJSONBody(w, r, &updates) {
		return
	}

	updated, err := h.ticketService.UpdateTicketInfo(r.Context(), uuid, updates)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *TicketHandler) GetReasonInfoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	reasonInfo, err := h.ticketService.GetReasonInfoByID(r.Context(), id)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, reasonInfo)
}

func (h *TicketHandler) GetTicketContactPerson(w http.ResponseWriter, r *http.Request) {
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

	contact, err := h.ticketService.GetTicketContactPerson(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, contact)
}

func (h *TicketHandler) GetTicketsByField(w http.ResponseWriter, r *http.Request) {
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

	fieldUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	ticketIDs, err := h.ticketService.GetTicketsByField(r.Context(), field, fieldUUID)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, ticketIDs)
}
