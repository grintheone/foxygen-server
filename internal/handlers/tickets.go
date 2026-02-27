package handlers

import (
	"encoding/json"
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
	limit, offset, sortByTitle, ok := parsePaginationParams(r, 100000)
	if !ok {
		clientError(w, http.StatusBadRequest)
		return
	}

	role, ok := middlewares.GetUserRoleFromContext(r.Context())
	if !ok {
		serverError(w, fmt.Errorf("Unable to check for user role"))
		return
	}

	executorID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		serverError(w, fmt.Errorf("Unable to check for user ID"))
		return
	}

	tickets, err := h.ticketService.ListAllTickets(r.Context(), executorID, role, limit, offset, sortByTitle)
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

	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		serverError(w, fmt.Errorf("no user ID present in context"))
	}

	var updates models.TicketUpdates
	updates.ID = uuid

	if !decodeJSONBody(w, r, &updates) {
		return
	}

	err = h.ticketService.UpdateTicketInfo(r.Context(), updates, userID)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, uuidStr)
}

func (h *TicketHandler) CloseTicket(w http.ResponseWriter, r *http.Request) {
	var ticketInfo models.CloseTicket

	if !decodeJSONBody(w, r, &ticketInfo) {
		return
	}

	uuidStr, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		serverError(w, fmt.Errorf("no user ID present in context"))
	}

	currentUserID, err := uuid.Parse(uuidStr)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	err = h.ticketService.CloseTicket(r.Context(), ticketInfo, currentUserID)
	if err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(200)
}

func (h *TicketHandler) GetTicketReasons(w http.ResponseWriter, r *http.Request) {
	reasons, err := h.ticketService.GetTicketReasons(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, reasons)
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

	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		serverError(w, fmt.Errorf("Unable to check for user ID"))
		return
	}

	f := r.URL.Query().Get("filters")

	var filters models.TicketFilters
	err = json.Unmarshal([]byte(f), &filters)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	tickets, err := h.ticketService.GetTicketsByField(r.Context(), field, fieldUUID, filters, userID)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tickets)
}
