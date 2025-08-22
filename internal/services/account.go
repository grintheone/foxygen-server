package services

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	accountRepo *repository.AccountRepository
}

func NewAccountService(ar *repository.AccountRepository) *AccountService {
	return &AccountService{accountRepo: ar}
}

// CreateUser is called from your HTTP handler.
// It takes the plain text password, hashes it, and defines the business rules for role assignment.
func (s *AccountService) CreateUser(ctx context.Context, username, plainPassword, database string, requestedRoleNames []string) (*models.Account, error) {
	// 1. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 2. Map role NAMES to role IDs based on your business logic.
	// This is better than letting the client send arbitrary role IDs.
	var roleIDsToAssign []int
	for _, roleName := range requestedRoleNames {
		switch roleName {
		case "admin":
			roleIDsToAssign = append(roleIDsToAssign, 1)
		case "coordinator":
			roleIDsToAssign = append(roleIDsToAssign, 2)
		case "user":
			roleIDsToAssign = append(roleIDsToAssign, 3)
		default:
			// Reject unknown roles
			return nil, errors.New("invalid role requested: " + roleName)
		}
	}

	// 3. Ensure every user has at least the 'user' role
	hasUserRole := slices.Contains(roleIDsToAssign, 3)
	if !hasUserRole {
		roleIDsToAssign = append(roleIDsToAssign, 3) // Assign the basic 'user' role
	}

	// 4. Create the account model and call the repository
	newAccount := &models.Account{
		Username:     username,
		Database:     database,
		PasswordHash: string(hashedPassword),
		Disabled:     false,
	}

	createdAccount, err := s.accountRepo.CreateAccountWithRoles(ctx, newAccount, roleIDsToAssign)
	if err != nil {
		return nil, err
	}

	// For security, don't return the password hash
	createdAccount.PasswordHash = ""
	return createdAccount, nil
}

// GetUserByUsername is the public method called by the AuthService.
// It handles the business logic for fetching a user.
func (s *AccountService) GetAccountByUsername(ctx context.Context, username string) (*models.Account, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	account, err := s.accountRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("service error fetching user: %w", err)
	}

	// The repository returns nil if the user is not found.
	// This is not an error, but a normal case we need to handle.
	if account == nil {
		return nil, nil
	}

	// You could add any additional business logic here.
	// For example, checking if the account is disabled before returning it?
	// Though the AuthService handles that, it might be better here.
	// if account.Disabled {
	// 	return nil, fmt.Errorf("account is disabled")
	// }

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
