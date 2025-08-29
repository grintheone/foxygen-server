package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type CommentService struct {
	commentRepo repository.CommentsRepository
}

func NewCommentService(r repository.CommentsRepository) *CommentService {
	return &CommentService{commentRepo: r}
}

func (s *CommentService) GetCommentByIds(ctx context.Context, ids string) (*[]models.Comment, error) {
	idsArr := strings.Split(ids, ",")

	// Convert each string to integer
	numberIds := make([]int, len(idsArr))
	for i, s := range idsArr {
		num, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			fmt.Printf("Error converting %s: %v\n", s, err)
			continue
		}
		numberIds[i] = num
	}

	comments, err := s.commentRepo.GetCommentByIds(ctx, numberIds)
	if err != nil {
		return nil, fmt.Errorf("service error fetching comments: %w", err)
	}

	return comments, nil
}

func (s *CommentService) NewComment(ctx context.Context, data models.Comment) (*models.Comment, error) {
	comment, err := s.commentRepo.NewComment(ctx, data)
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

	err = s.commentRepo.DeleteComment(ctx, numberId)
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

	err = s.commentRepo.UpdateComment(ctx, numberId, payload)
	if err != nil {
		return fmt.Errorf("service error updating the comment: %w", err)
	}

	return nil
}
