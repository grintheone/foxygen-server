package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/services"
)

type DepartmentHandler struct {
	service *services.DepartmentService
}

func (h *DepartmentHandler) ListAllDepartments(w http.ResponseWriter, r *http.Request) {
	departments, err := h.service.ListAllDepartments(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, departments)
}

func (h *DepartmentHandler) GetDepartmentByID(w http.ResponseWriter, r *http.Request) {
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

	department, err := h.service.GetDepartmentByID(r.Context(), uuid)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, department)
}
