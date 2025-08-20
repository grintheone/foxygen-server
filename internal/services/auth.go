package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Define token expiration times
const (
	AccessTokenExpiry  = 30 * time.Minute
	RefreshTokenExpiry = 14 * 24 * time.Hour // 14 days
)

// LoginResponse defines the structure of a successful login response
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"` // Usually "Bearer"
	ExpiresIn    int64  `json:"expiresIn"` // Seconds until access token expiry
}

var ErrInvalidCredentials = errors.New("invalid username or password")

type AuthService struct {
	accountService *AccountService
	jwtSecret      []byte
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthService(as *AccountService, jwtSecret string) *AuthService {
	return &AuthService{
		accountService: as,
		jwtSecret:      []byte(jwtSecret),
	}
}

// generateAccessToken creates a short-lived token for API access
func (s *AuthService) generateAccessToken(user *models.Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.UserID.String(),
		"exp":      time.Now().Add(AccessTokenExpiry).Unix(),
		"username": user.Username,
		"roles":    user.GetRoleNames(),
		"type":     "access", // Explicitly mark token type
	})
	return token.SignedString(s.jwtSecret)
}

// generateRefreshToken creates a long-lived token only used for getting new access tokens
func (s *AuthService) generateRefreshToken(user *models.Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.UserID.String(),
		"exp":  time.Now().Add(RefreshTokenExpiry).Unix(),
		"type": "refresh", // Explicitly mark token type
		// Note: Refresh tokens should contain minimal claims for security
	})
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	// 1. Find the user by username.
	// You will need to create a `GetUserByUsername` method in your AccountService/Repository.
	user, err := s.accountService.GetAccountByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	if user == nil {
		// For security, hash a dummy password to prevent timing attacks
		// that reveal whether a username exists based on response time.
		dummyHash := "$2a$10$dummyhashdummyhashdummyhashdummyhashdummyhashdummyha"
		bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(password))
		return nil, ErrInvalidCredentials
	}

	// 2. Check if the account is disabled.
	if user.Disabled {
		return nil, errors.New("account is disabled")
	}

	// 3. Compare the provided password with the stored hash.
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate both access and refresh tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	response := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(AccessTokenExpiry.Seconds()),
	}

	return response, nil
}

// ValidateToken is used by your middleware to protect routes.
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// RefreshAccessToken validates a refresh token and issues a new access token
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshTokenString string) (*LoginResponse, error) {
	// Validate the refresh token
	token, err := s.validateTokenAndType(refreshTokenString, "refresh")
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Extract the user ID from the refresh token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("user ID not found in token")
	}

	// Now parse the string into a UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format in token: %w", err)
	}

	// Fetch the user from the database to ensure they still exist and are active
	user, err := s.accountService.GetUserByID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user no longer exists")
	}
	if user.Disabled {
		return nil, errors.New("account is disabled")
	}

	// Generate a new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	response := &LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(AccessTokenExpiry.Seconds()),
		// Return the same refresh token - it remains valid until its own expiration
		RefreshToken: refreshTokenString,
	}

	return response, nil
}

// validateTokenAndType validates a token and checks its type claim
func (s *AuthService) validateTokenAndType(tokenString, expectedType string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check token type
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	return token, nil
}

// ValidateAccessToken is a wrapper for validating access tokens specifically
func (s *AuthService) ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	return s.validateTokenAndType(tokenString, "access")
}
