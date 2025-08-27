package repository

import (
	"context"
	"time"

	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CommentsRepository interface {
	GetCommentByIds(ctx context.Context, ids []int) (*[]models.Comment, error)
	NewComment(ctx context.Context, data models.Comment) (*models.Comment, error)
	DeleteComment(ctx context.Context, id int) error
	UpdateComment(ctx context.Context, updates models.CommentUpdate) error
}

type commentsRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) CommentsRepository {
	return &commentsRepository{db}
}

func (r *commentsRepository) GetCommentByIds(ctx context.Context, ids []int) (*[]models.Comment, error) {
	query := `
        SELECT id, author_id, reference_id, text, created_at
        FROM comments
        WHERE id = ANY($1)
        ORDER BY created_at DESC
    `
	var comments []models.Comment
	err := r.db.SelectContext(ctx, &comments, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	return &comments, nil
}

func (r *commentsRepository) NewComment(ctx context.Context, data models.Comment) (*models.Comment, error) {
	query := `
        INSERT INTO comments (author_id, reference_id, text, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id, author_id, reference_id, text, created_at
    `

	var comment models.Comment

	err := r.db.Get(&comment, query, data.AuthorID, data.ReferenceID, data.Text, time.Now())
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *commentsRepository) DeleteComment(ctx context.Context, id int) error {
	query := `DELETE FROM comments WHERE id = $1;`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *commentsRepository) UpdateComment(ctx context.Context, updates models.CommentUpdate) error {
	query := `UPDATE comments SET text = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, updates.Text, updates.ID)
	if err != nil {
		return err
	}

	return nil
}
