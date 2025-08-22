package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/grintheone/foxygen-server/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	response, err := h.authService.Authorize(r.Context(), credentials.Username, credentials.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			http.Error(w, "Login failed", http.StatusUnauthorized)
		} else {
			http.Error(w, "Login failed", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Refresh handles token refresh requests
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if request.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	response, err := h.authService.RefreshAccessToken(r.Context(), request.RefreshToken)
	if err != nil {
		http.Error(w, "Token refresh failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
