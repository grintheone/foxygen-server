package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grintheone/foxygen-server/internal/middlewares"
	"github.com/grintheone/foxygen-server/internal/services"
)

func NewRouter(accountService *services.AccountService, authService *services.AuthService) http.Handler {
	r := chi.NewRouter()
	// Initialize handlers
	accountHandler := &AccountHandler{accountService: accountService}
	authHandler := &AuthHandler{authService: authService}

	// Global middleware (applied to all routes)
	r.Use(middleware.Logger)    // Logs incoming requests
	r.Use(middleware.Recoverer) // Recovers from panics
	r.Use(middleware.RequestID) // Adds a request ID to each request

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("tls working"))
		})
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login) // Main handler for further operations with the app
		r.Post("/refreshToken", authHandler.Refresh)
	})

	// Router that requires authentication
	r.Route("/app", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(authService))

		// Main app handlers
	})

	// Protected router for account specific operations
	r.Route("/account", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(authService))

		r.Patch("/change-password", accountHandler.ChangePassword)
	})

	// Router that requires authentication and admin role
	r.Route("/admin", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(authService))
		r.Use(middlewares.RequireRole("admin"))

		r.Post("/create-account", accountHandler.CreateAccount)
		r.Patch("/change-account-status", accountHandler.ChangeAccountStatus)
	})

	return r
}
