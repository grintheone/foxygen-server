package models

import "github.com/google/uuid"

type Department struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Title string    `json:"title" db:"title"`
}
