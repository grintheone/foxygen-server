package handlers

import (
	"fmt"
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

	r.Group(func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/register", accountHandler.CreateAccount) // Or make this admin-only
	})

	// Protected routes (authentication required)
	r.Group(func(r chi.Router) {
		// Apply auth middleware to this group
		r.Use(middlewares.AuthMiddleware(authService))

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			fmt.Print(r.Context())
			w.Write([]byte("cool"))
		})

		// r.Get("/profile", accountHandler.GetProfile)
		// r.Put("/profile", accountHandler.UpdateProfile)

		// Admin-only routes (authentication + role check)
		r.Group(func(r chi.Router) {
			// Apply role-based middleware on top of auth
			r.Use(middlewares.RequireRole("admin"))

			// r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
			// 	fmt.Print(r.Context())
			// 	w.Write([]byte("admin"))
			// })

			// r.Get("/admin/users", accountHandler.ListAllUsers)
			// r.Post("/admin/users", accountHandler.CreateUser) // Admin creates users
			// r.Put("/admin/users/{userID}", accountHandler.UpdateUser)
		})
	})

	return r
}
