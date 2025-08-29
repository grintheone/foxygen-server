package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/services"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func (h *CommentHandler) GetCommentByIds(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("ids")
	if ids == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	comments, err := h.commentService.GetCommentByIds(r.Context(), ids)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, comments)
}

func (h *CommentHandler) NewComment(w http.ResponseWriter, r *http.Request) {
	var request models.Comment

	if !decodeJSONBody(w, r, &request) {
		return
	}

	createdComment, err := h.commentService.NewComment(r.Context(), request)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createdComment)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	err := h.commentService.DeleteComment(r.Context(), id)
	if err != nil {
		serverError(w, err)
		return
	}

	var response struct {
		ID string `json:"id"`
	}

	response.ID = id

	writeJSON(w, http.StatusOK, response)
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")

	if ID == "" {
		clientError(w, http.StatusBadRequest)
		return
	}

	var payload models.CommentUpdate

	if !decodeJSONBody(w, r, &payload) {
		return
	}

	err := h.commentService.UpdateComment(r.Context(), ID, payload)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, payload)
}
