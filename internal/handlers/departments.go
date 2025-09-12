package handlers

import (
	"net/http"

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
