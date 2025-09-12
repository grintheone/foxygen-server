package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/services"
)

type AttachmentHandler struct {
	attachmentService *services.AttachmentService
	uploadDir         string
}

func (h *AttachmentHandler) LoadImageByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	attachment, err := h.attachmentService.GetAttachmentByID(r.Context(), id)
	if err != nil {
		serverError(w, err)
		return
	}

	writeFile(w, r, *attachment)
}

func (h *AttachmentHandler) GetAttachmentsByRefID(w http.ResponseWriter, r *http.Request) {
	refID := chi.URLParam(r, "refID")
	if refID == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	refUUID, err := uuid.Parse(refID)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	attachments, err := h.attachmentService.GetAttachmentsByRefID(r.Context(), refUUID)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, attachments)
}

func (h *AttachmentHandler) UploadMultiple(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with larger max memory (e.g., 50MB for multiple files)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		log.Print(err)
		clientError(w, http.StatusBadRequest)
		return
	}

	// Get all files from the form
	form := r.MultipartForm
	files := form.File["attachments"] // Use the same field name as frontend

	if len(files) == 0 {
		log.Print("No files in request")
		clientError(w, http.StatusBadRequest)
		return
	}

	// Get ref_id from form data
	refID := r.FormValue("ref_id")
	if refID == "" {
		log.Print("Ref ID not present")
		clientError(w, http.StatusBadRequest)
		return
	}

	refUUID, err := uuid.Parse(refID)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	// Upload all files
	attachments, err := h.attachmentService.UploadMultipleFiles(r.Context(), files, h.uploadDir, refUUID)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, attachments)
}

func (h *AttachmentHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 10MB max memory
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Print(err)
		clientError(w, http.StatusBadRequest)
		return
	}

	// Get ref_id from form data
	refID := r.FormValue("ref_id")
	if refID == "" {
		log.Print("Ref ID not present")
		clientError(w, http.StatusBadRequest)
		return
	}

	refUUID, err := uuid.Parse(refID)
	if err != nil {
		clientError(w, http.StatusBadRequest)
		return
	}

	// Get the file from form data
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Print(err)
		clientError(w, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload file
	attachment, err := h.attachmentService.UploadFile(r.Context(), fileHeader, h.uploadDir, refUUID)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, attachment)
}
