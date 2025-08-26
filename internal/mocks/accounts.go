package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) CreateAccountWithRoles(ctx context.Context, account *models.Account, roleID int) (*models.Account, error) {
	args := m.Called(ctx, account, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByUsername(ctx context.Context, username string) (*models.Account, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.Account, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) ChangeAccountPassword(ctx context.Context, userID uuid.UUID, hash string) error {
	args := m.Called(ctx, userID, hash)
	return args.Error(0)
} 

func (m *MockAccountRepository) ChangeAccountStatus(ctx context.Context, userID uuid.UUID, disabled bool) error {
	args := m.Called(ctx, userID, disabled)
	return args.Error(0)
}
