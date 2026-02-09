package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/grintheone/foxygen-server/internal/config"
	"github.com/grintheone/foxygen-server/internal/handlers"
	"github.com/grintheone/foxygen-server/internal/repository"
	"github.com/grintheone/foxygen-server/internal/services"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type App struct {
	Router     http.Handler
	DB         *sqlx.DB
	ImportFile *string
}

func NewApp(cfg *config.Config, importFile *string) (*App, error) {
	db, err := sqlx.Open("postgres", cfg.Database.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetConnMaxIdleTime(5 * time.Minute)

	accountRepo := repository.NewAccountRepository(db)
	accountService := services.NewAccountService(accountRepo)
	authService := services.NewAuthService(accountService, cfg.Server.Secret)

	// User
	userRepo := repository.NewUsersRepository(db)
	userService := services.NewUserService(userRepo)

	// Client
	clientRepo := repository.NewClientRepository(db)
	clientService := services.NewClientService(clientRepo)

	// Comment
	commentRepo := repository.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo)

	// Contact
	contactRepo := repository.NewContactRepository(db)
	contactService := services.NewContactService(contactRepo)

	// Devices
	deviceRepo := repository.NewDeviceRepository(db)
	deviceService := services.NewDeviceService(deviceRepo)

	// Classificators
	classificatorRepo := repository.NewClassificatorRepository(db)
	classificatorService := services.NewClassificatorService(classificatorRepo)

	// Tickets
	ticketRepo := repository.NewTicketRepository(db)
	ticketService := services.NewTicketService(ticketRepo)

	// Attachments
	attachmentsRepo := repository.NewAttachmentRepository(db)
	attachmentService := services.NewAttachmentService(attachmentsRepo)

	// Departments
	departmentRepo := repository.NewDepartmentRepo(db)
	departmentService := services.NewDepartmentService(departmentRepo)

	// Agreements
	agreementRepo := repository.NewAgreementRepo(db)
	agreementService := services.NewAgreementService(agreementRepo)

	// Regions
	regionsRepo := repository.NewRegionRepo(db)
	// Researh Type
	researchTypeRepo := repository.NewResearchTypeRepo(db)
	// Manufacturer
	manufacturerRepo := repository.NewManufacturerRepo(db)

	importService := services.NewImportService(
		*departmentService,
		*userService,
		*accountService,
		*clientService,
		*contactService,
		*classificatorService,
		*deviceService,
		ticketRepo,
		regionsRepo,
		researchTypeRepo,
		manufacturerRepo,
	)
	// Import data if flag is provided
	if *importFile != "" {
		log.Printf("Importing data from: %s", *importFile)
		if err := importService.ImportFromFile(*importFile); err != nil {
			log.Fatalf("Failed to import data: %v", err)
		}
		log.Println("Data import completed successfully")
	}

	r := handlers.NewRouter(
		accountService,
		authService,
		userService,
		clientService,
		commentService,
		contactService,
		deviceService,
		classificatorService,
		ticketService,
		attachmentService,
		departmentService,
		agreementService,
	)

	return &App{Router: r, DB: db}, nil
}

func (a *App) Close() error {
	return a.DB.Close()
}
