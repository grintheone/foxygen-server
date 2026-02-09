package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/middlewares"
	"github.com/grintheone/foxygen-server/internal/services"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
		Role     string `json:"role"`
	}

	if !decodeJSONBody(w, r, &request) {
		return
	}

	newUser, err := h.accountService.CreateUser(
		r.Context(),
		request.Username,
		request.Password,
		request.Role,
		nil,
	)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newUser)
}

func (h *AccountHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var request struct {
		New string `json:"new"`
		Old string `json:"old"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unable to prove your identity", http.StatusForbidden)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusForbidden)
		return
	}

	err = h.accountService.ChangeAccountPassword(r.Context(), userUUID, request.New, request.Old)
	if err != nil {
		log.Printf("handler: %v", err)
		http.Error(w, "Unable to change password", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AccountHandler) ChangeAccountStatus(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Disabled bool `json:"disabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unable to prove your identity", http.StatusForbidden)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusForbidden)
		return
	}

	err = h.accountService.ChangeAccountStatus(r.Context(), userUUID, request.Disabled)
	if err != nil {
		log.Printf("handler: %v", err)
		http.Error(w, "Unable to change account status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
