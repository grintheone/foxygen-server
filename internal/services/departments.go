package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
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

func (s *DepartmentService) GetDepartmentByID(ctx context.Context, uuid uuid.UUID) (*models.Department, error) {
	department, err := s.repo.GetDepartmentByID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("service error fetching a department: %w", err)
	}

	return department, nil
}

func (s *DepartmentService) AddNewDepartment(ctx context.Context, data models.Department) error {
	err := s.repo.AddNewDepartment(ctx, data)
	if err != nil {
		return fmt.Errorf("service error adding a new department: %w", err)
	}

	return nil
}
