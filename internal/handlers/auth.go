package handlers

import (
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

	if !decodeJSONBody(w, r, &credentials) {
		return
	}

	response, err := h.authService.Authorize(r.Context(), credentials.Username, credentials.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			clientError(w, http.StatusUnauthorized)
		} else {
			serverError(w, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}

// Refresh handles token refresh requests
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refreshToken"`
	}

	if !decodeJSONBody(w, r, &request) {
		return
	}

	if request.RefreshToken == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	response, err := h.authService.RefreshAccessToken(r.Context(), request.RefreshToken)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}
