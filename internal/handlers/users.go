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

type UserHandler struct {
	userService *services.UserService
}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "userID")

	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	profile, err := h.userService.GetUserProfile(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

func (h *UserHandler) ListDepartmentUsers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		serverError(w, fmt.Errorf("No user ID is present in context"))
	}

	users, err := h.userService.ListDepartmentUsers(r.Context(), userID)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, users)
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
