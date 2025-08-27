package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		serverError(w, err)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userUUID)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListUsers(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		serverError(w, err)
		return
	}

	err = h.userService.DeleteUser(r.Context(), userUUID)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, userID)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		serverError(w, err)
		return
	}

	var user models.User
	if !decodeJSONBody(w, r, &user) {
		return
	}

	user.UserID = userUUID

	err = h.userService.UpdateUser(r.Context(), user)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
