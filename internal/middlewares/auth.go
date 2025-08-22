package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grintheone/foxygen-server/internal/services"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// UserIDKey is the key for storing user ID in the request context.
	UserIDKey contextKey = "user_id"
	// UsernameKey is the key for storing username in the request context.
	UsernameKey contextKey = "username"
	// UserRolesKey is the key for storing user roles in the request context.
	UserRoleKey contextKey = "user_role"
)

func handleAuthError(w http.ResponseWriter, err error) {
	log.Print(err)
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

// AuthMiddleware creates a middleware that validates JWT tokens.
func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handleAuthError(w, errors.New("auth header is empty"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				handleAuthError(w, errors.New("invalid token type"))
				return
			}

			tokenString := parts[1]

			token, err := authService.ValidateAccessToken(tokenString)
			if err != nil {
				handleAuthError(w, err)
				return
			}

			// Extract claims and add them to the request context
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx := r.Context()
				if sub, exists := claims["sub"]; exists {
					ctx = context.WithValue(ctx, UserIDKey, sub)
				}
				if username, exists := claims["username"]; exists {
					ctx = context.WithValue(ctx, UsernameKey, username)
				}
				if role, exists := claims["role"]; exists {
					ctx = context.WithValue(ctx, UserRoleKey, role)
				}
				r = r.WithContext(ctx)
			} else {
				handleAuthError(w, errors.New("wrong claims"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole is an example of additional middleware for role-based access control (RBAC)
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := GetUserRoleFromContext(r.Context())
			if !ok {
				log.Print("Access denied: no role information")
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			if role != requiredRole {
				log.Print("Access denied: insufficient permissions")
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Helper functions to extract data from the context in your handlers
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
}

func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}

func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	roles, ok := ctx.Value(UserRoleKey).(string)
	return roles, ok
}
