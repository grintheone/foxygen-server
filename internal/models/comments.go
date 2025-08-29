package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID          int       `json:"id,omitzero" db:"id"`
	AuthorID    uuid.UUID `json:"author_id" db:"author_id"`
	ReferenceID uuid.UUID `json:"reference_id" db:"reference_id"`
	Text        string    `json:"text" db:"text"`
	CreatedAt   time.Time `json:"created_at,omitzero" db:"created_at"`
}

type CommentUpdate struct {
	Text string `json:"text" db:"text"`
}
