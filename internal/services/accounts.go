package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(ar repository.AccountRepository) *AccountService {
	return &AccountService{accountRepo: ar}
}

func (s *AccountService) CreateUser(ctx context.Context, username, password, database string, role string) (*models.Account, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var roleID int

	switch role {
	case "admin":
		roleID = 1
	case "coordinator":
		roleID = 2
	case "user":
		roleID = 3
	default:
		return nil, errors.New("invalid role requested: " + role)
	}

	newAccount := &models.Account{
		Username:     username,
		Database:     database,
		PasswordHash: string(hashedPassword),
	}

	createdAccount, err := s.accountRepo.CreateAccountWithRoles(ctx, newAccount, roleID)
	if err != nil {
		return nil, err
	}

	// For security, don't return the password hash
	createdAccount.PasswordHash = ""
	return createdAccount, nil
}

func (s *AccountService) GetAccountByUsername(ctx context.Context, username string) (*models.Account, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	account, err := s.accountRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("service error fetching user: %w", err)
	}

	if account == nil {
		return nil, nil
	}

	if account.Disabled {
		return nil, fmt.Errorf("account is disabled")
	}

	return account, nil
}

func (s *AccountService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service error fetching user by ID: %w", err)
	}
	if account == nil {
		return nil, nil
	}

	if account.Disabled {
		return nil, fmt.Errorf("account is disabled")
	}

	return account, nil
}

func (s *AccountService) ChangeAccountPassword(ctx context.Context, userID uuid.UUID, new, old string) error {
	account, err := s.accountRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("service error fetching user by ID: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(old))
	if err != nil {
		return ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(new), bcrypt.DefaultCost)

	err = s.accountRepo.ChangeAccountPassword(ctx, userID, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountService) ChangeAccountStatus(ctx context.Context, userID uuid.UUID, disabled bool) error {
	err := s.accountRepo.ChangeAccountStatus(ctx, userID, disabled)
	if err != nil {
		return err
	}

	return nil
}
