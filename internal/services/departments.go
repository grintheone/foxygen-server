package services

import (
	"context"
	"fmt"

	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type DepartmentService struct {
	repo repository.DepartmentRepo
}

func NewDepartmentService(repo repository.DepartmentRepo) *DepartmentService {
	return &DepartmentService{repo}
}

func (s *DepartmentService) ListAllDepartments(ctx context.Context) ([]*models.Department, error) {
	departments, err := s.repo.ListAllDepartments(ctx)
	if err != nil {
		return nil, fmt.Errorf("service error fetching departments: %w", err)
	}

	return departments, nil
}
