package models

import "github.com/google/uuid"

type Attachment struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	MediaType string    `json:"mediaType" db:"media_type"`
	Ext       string    `json:"ext" db:"ext"`
	RefID     uuid.UUID `json:"refID" db:"ref_id"`
}
