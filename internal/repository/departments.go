package repository

import (
	"context"

	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type DepartmentRepo interface {
	ListAllDepartments(ctx context.Context) ([]*models.Department, error)
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
