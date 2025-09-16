package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type DepartmentRepo interface {
	ListAllDepartments(ctx context.Context) ([]*models.Department, error)
	GetDepartmentByID(ctx context.Context, uuid uuid.UUID) (*models.Department, error)
}

type departmentRepo struct {
	db *sqlx.DB
}

func NewDepartmentRepo(db *sqlx.DB) *departmentRepo {
	return &departmentRepo{db}
}

func (r *departmentRepo) ListAllDepartments(ctx context.Context) ([]*models.Department, error) {
	var departments []*models.Department

	err := r.db.SelectContext(ctx, &departments, `SELECT * FROM departments`)
	if err != nil {
		return nil, err
	}

	return departments, nil
}

func (r *departmentRepo) GetDepartmentByID(ctx context.Context, uuid uuid.UUID) (*models.Department, error) {
	var department models.Department

	fmt.Print(uuid)

	err := r.db.GetContext(ctx, &department, `SELECT * FROM departments WHERE id = $1`, uuid)
	if err != nil {
		return nil, err
	}

	return &department, nil
}
