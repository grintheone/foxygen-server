package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grintheone/foxygen-server/internal/middlewares"
	"github.com/grintheone/foxygen-server/internal/services"
)

func NewRouter(
	accountService *services.AccountService,
	authService *services.AuthService,
	userService *services.UserService,
	clientService *services.ClientService,
	commentService *services.CommentService,
	contactService *services.ContactService,
) http.Handler {
	r := chi.NewRouter()
	// Initialize handlers
	accountHandler := &AccountHandler{accountService}
	authHandler := &AuthHandler{authService}
	userHandler := &UserHandler{userService}
	clientHandler := &ClientHandler{clientService}
	commentHandler := &CommentHandler{commentService}
	contactHandler := &ContactHandler{contactService}

	// Global middleware (applied to all routes)
	r.Use(middleware.Logger)    // Logs incoming requests
	r.Use(middleware.Recoverer) // Recovers from panics
	r.Use(middleware.RequestID) // Adds a request ID to each request

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login) // Main handler for further operations with the app
		r.Post("/refreshToken", authHandler.Refresh)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(authService))

		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.ListUsers)
				r.Get("/{userID}", userHandler.GetByID)
				r.Delete("/{userID}", userHandler.DeleteUser)
				r.Patch("/{userID}", userHandler.UpdateUser)
			})

			r.Route("/clients", func(r chi.Router) {
				r.Get("/", clientHandler.ListClients)
			})

			r.Route("/comments", func(r chi.Router) {
				r.Get("/", commentHandler.GetCommentByIds)
				r.Post("/", commentHandler.NewComment)
				r.Patch("/", commentHandler.UpdateComment)
				r.Delete("/{id}", commentHandler.DeleteComment)
			})

			r.Route("/contacts", func(r chi.Router) {
				r.Get("/{clientID}", contactHandler.GetAllByClientID)
				r.Post("/", contactHandler.CreateContact)
				r.Delete("/{id}", contactHandler.DeleteContact)
				r.Patch("/{id}", contactHandler.UpdateContact)
			})

			r.Route("/accounts", func(r chi.Router) {
				r.Patch("/password", accountHandler.ChangePassword)
			})

			// Router that requires authentication and admin role
			r.Route("/admin", func(r chi.Router) {
				r.Use(middlewares.RequireRole("admin"))

				r.Route("/accounts", func(r chi.Router) {
					r.Post("/", accountHandler.CreateAccount)
					r.Patch("/status", accountHandler.ChangeAccountStatus)
				})
			})
		})
	})

	return r
}
