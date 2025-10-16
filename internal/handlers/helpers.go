package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/grintheone/foxygen-server/internal/models"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func serverError(w http.ResponseWriter, err error) {
	// trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// log.Print(trace)
	log.Print(err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user.
func clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func notFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func decodeJSONBody[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		log.Print(err)
		clientError(w, http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()
	return true
}

func writeJSON[T any](w http.ResponseWriter, status int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		serverError(w, err)
	}
}

func writeFile(w http.ResponseWriter, r *http.Request, data models.Attachment) {
	// Set headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename="+data.OriginalName)
	w.Header().Set("Content-Type", data.MimeType)
	w.Header().Set("Content-Length", strconv.FormatInt(data.FileSize, 10))

	http.ServeFile(w, r, data.FilePath)
}
