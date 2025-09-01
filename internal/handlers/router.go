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
	deviceService *services.DeviceService,
	classificatorService *services.ClassificatorService,
) http.Handler {
	r := chi.NewRouter()
	// Initialize handlers
	accountHandler := &AccountHandler{accountService}
	authHandler := &AuthHandler{authService}
	userHandler := &UserHandler{userService}
	clientHandler := &ClientHandler{clientService}
	commentHandler := &CommentHandler{commentService}
	contactHandler := &ContactHandler{contactService}
	deviceHandler := &DeviceHandler{deviceService}
	classificatorHandler := &ClassificatorHandler{classificatorService}

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
				r.Post("/", clientHandler.CreateClient)
				r.Patch("/{uuid}", clientHandler.UpdateClient)
				r.Delete("/{uuid}", clientHandler.DeleteClient)
			})

			r.Route("/comments", func(r chi.Router) {
				r.Get("/{uuid}", commentHandler.GetCommentsByReferenceID)
				r.Post("/", commentHandler.NewComment)
				r.Patch("/{id}", commentHandler.UpdateComment)
				r.Delete("/{id}", commentHandler.DeleteComment)
			})

			r.Route("/devices", func(r chi.Router) {
				r.Get("/", deviceHandler.GetAllDevices)
				r.Post("/", deviceHandler.CreateNewDevice)
				r.Get("/{uuid}", deviceHandler.GetDeviceByID)
				r.Delete("/{uuid}", deviceHandler.RemoveDeviceByID)
				r.Patch("/{uuid}", deviceHandler.UpdateDeviceByID)
			})

			r.Route("/classificators", func(r chi.Router) {
				r.Get("/{uuid}", classificatorHandler.GetClassificatorByID)
				r.Get("/devices/{uuid}", classificatorHandler.GetDevicesByClassificatorID)
				r.Post("/", classificatorHandler.NewClassificator)
				r.Delete("/{uuid}", classificatorHandler.RemoveClassificatorByID)
				r.Patch("/{uuid}", classificatorHandler.UpdateClassificatorInfo)
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
