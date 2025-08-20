package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/grintheone/foxygen-server/internal/services"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string   `json:"username"`
		Password string   `json:"password"`
		Database string   `json:"database"`
		Roles    []string `json:"roles"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Call the service layer
	newUser, err := h.accountService.CreateUser(r.Context(), request.Username, request.Password, request.Database, request.Roles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
