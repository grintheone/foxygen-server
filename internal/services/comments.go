package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type CommentService struct {
	repo repository.CommentsRepository
}

func NewCommentService(r repository.CommentsRepository) *CommentService {
	return &CommentService{repo: r}
}

func (s *CommentService) GetCommentsByReferenceID(ctx context.Context, uuid uuid.UUID) (*[]models.Comment, error) {
	comments, err := s.repo.GetCommentsByReferenceID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("service error fetching comments: %w", err)
	}

	return comments, nil
}

func (s *CommentService) NewComment(ctx context.Context, data models.Comment) (*models.Comment, error) {
	comment, err := s.repo.NewComment(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("service error creating a comment: %w", err)
	}

	return comment, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, id string) error {
	numberId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("can't convert id to number: %w", err)
	}

	err = s.repo.DeleteComment(ctx, numberId)
	if err != nil {
		return fmt.Errorf("service error deleting a comment: %w", err)
	}

	return nil
}

func (s *CommentService) UpdateComment(ctx context.Context, id string, payload models.CommentUpdate) error {
	numberId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("can't convert id to number: %w", err)
	}

	err = s.repo.UpdateComment(ctx, numberId, payload)
	if err != nil {
		return fmt.Errorf("service error updating the comment: %w", err)
	}

	return nil
}
