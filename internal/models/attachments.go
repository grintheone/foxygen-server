package models

import "github.com/google/uuid"

type Attachment struct {
	ID           int       `json:"id" db:"id"`
	FileName     string    `json:"file_name" db:"file_name"`
	OriginalName string    `json:"original_name" db:"original_name"`
	FileSize     int64     `json:"file_size" db:"file_size"`
	FilePath     string    `json:"file_path" db:"file_path"`
	MimeType     string    `json:"mime_type" db:"mime_type"`
	RefID        uuid.UUID `json:"ref_id" db:"ref_id"`
}
