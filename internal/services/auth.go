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
	AccessTokenExpiry = 1 * time.Minute
	// AccessTokenExpiry  = 30 * time.Hour
	RefreshTokenExpiry = 10 * 24 * time.Hour // 10 days
)

// LoginResponse defines the structure of a successful login response
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"` // Usually "Bearer"
	ExpiresIn    int64  `json:"expiresIn"` // Seconds until access token expiry
}

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidToken       = errors.New("invalid token")
)

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
		"role":     user.Role,
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

func (s *AuthService) Authorize(ctx context.Context, username, password string) (*LoginResponse, error) {
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

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

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

// RefreshAccessToken validates a refresh token and issues a new access token
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshTokenString string) (*LoginResponse, error) {
	token, err := s.ValidateRefreshToken(refreshTokenString)
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

	// Generate a new access token
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

// validateTokenAndType validates a token and checks its type claim
func (s *AuthService) validateTokenAndType(tokenString, expectedType string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
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

func (s *AuthService) ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	return s.validateTokenAndType(tokenString, "refresh")
}
