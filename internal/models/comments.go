package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID          int       `json:"id" db:"id"`
	AuthorID    uuid.UUID `json:"author_id" db:"author_id"`
	ReferenceID uuid.UUID `json:"reference_id" db:"reference_id"`
	Text        string    `json:"text" db:"text"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CommentUpdate struct {
	ID   int    `json:"id" db:"id"`
	Text string `json:"text" db:"text"`
}
