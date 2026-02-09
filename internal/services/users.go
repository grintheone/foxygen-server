package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type UserService struct {
	userRepo repository.UsersRepository
}

func NewUserService(ur repository.UsersRepository) *UserService {
	return &UserService{userRepo: ur}
}

func (s *UserService) CreateNewUser(ctx context.Context, userData models.User) error {
	err := s.userRepo.CreateUser(ctx, userData)
	if err != nil {
		return fmt.Errorf("unable to create new user: %w", err)
	}

	return nil
}

func (s *UserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	profile, err := s.userRepo.GetProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service error fetching user profile: %w", err)
	}

	return profile, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service error fetching user by ID: %w", err)
	}
	if user == nil {
		return nil, nil
	}

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context) (*[]models.User, error) {
	users, err := s.userRepo.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("service error fetching all users: %w", err)
	}
	return users, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("service error when deleting user: %w", err)
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, user models.User) error {
	err := s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("service error updating the user: %w", err)
	}

	return nil
}

func (s *UserService) ListDepartmentUsers(ctx context.Context, userID string) ([]*models.User, error) {
	users, err := s.userRepo.ListDepartmentUsers(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service error retrieving the list of users: %w", err)
	}

	return users, nil
}
