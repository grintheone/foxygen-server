package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type ImportService struct {
	departmentService    DepartmentService
	userService          UserService
	accountService       AccountService
	clientService        ClientService
	contactService       ContactService
	classificatorService ClassificatorService
	deviceService        DeviceService
	ticketRepo           repository.TicketsRepository
	regionsRepo          repository.RegionsRepo
	researchTypeRepo     repository.ResearchTypeRepo
	manufacturerRepo     repository.ManufacturerRepo
}

func NewImportService(
	departmentService DepartmentService,
	userService UserService,
	accountService AccountService,
	clientService ClientService,
	contactService ContactService,
	classificatorService ClassificatorService,
	deviceService DeviceService,
	ticketRepo repository.TicketsRepository,
	regionsRepo repository.RegionsRepo,
	researchTypeRepo repository.ResearchTypeRepo,
	manufacturerRepo repository.ManufacturerRepo,
) *ImportService {
	return &ImportService{
		departmentService,
		userService,
		accountService,
		clientService,
		contactService,
		classificatorService,
		deviceService,
		ticketRepo,
		regionsRepo,
		researchTypeRepo,
		manufacturerRepo,
	}
}

// CouchDB export structure
type couchDBExport struct {
	TotalRows int          `json:"total_rows"`
	Offset    int          `json:"offset"`
	Rows      []couchDBRow `json:"rows"`
}

type couchDBRow struct {
	ID  string          `json:"id"`
	Key string          `json:"key"`
	Doc json.RawMessage `json:"doc"`
}

func extractAndParseUUID(rawID string) (uuid.UUID, error) {
	parts := strings.Split(rawID, "_1_")
	if len(parts) < 2 {
		return uuid.Nil, fmt.Errorf("invalid ID format")
	}

	return uuid.Parse(parts[1])
}

func (s *ImportService) processRowsInOrder(rows []couchDBRow) {
	// Separate rows by type first
	var departments, users, regions, researchType, manufacturers, clients, contacts, classificators, devices, tickets []couchDBRow

	for _, row := range rows {
		switch {
		case strings.HasPrefix(row.ID, "department_1_"):
			departments = append(departments, row)
		case strings.HasPrefix(row.ID, "user_1_"):
			users = append(users, row)
		case strings.HasPrefix(row.ID, "region_1_"):
			regions = append(regions, row)
		case strings.HasPrefix(row.ID, "researchType_1_"):
			researchType = append(researchType, row)
		case strings.HasPrefix(row.ID, "manufacturer_1_"):
			manufacturers = append(manufacturers, row)
		case strings.HasPrefix(row.ID, "client_1_"):
			clients = append(clients, row)
		case strings.HasPrefix(row.ID, "contact_1_"):
			contacts = append(contacts, row)
		case strings.HasPrefix(row.ID, "classificator_1_"):
			classificators = append(classificators, row)
		case strings.HasPrefix(row.ID, "device_1_"):
			devices = append(devices, row)
		case strings.HasPrefix(row.ID, "ticket_1_"):
			tickets = append(tickets, row)
		}
	}

	processed := 0
	failed := 0

	for _, row := range departments {
		if err := s.processDepartment(row.Doc); err != nil {
			failed++
			log.Printf("Failed to process row %s: %v", row.ID, err)
			continue
		}
		processed++
	}

	for i, row := range users {
		if err := s.processUser(row.Doc, i); err != nil {
			failed++
			log.Printf("Failed to process row %s: %v", row.ID, err)
			continue
		}
		processed++
	}

	for _, row := range regions {
		if err := s.processRegion(row.Doc); err != nil {
			failed++
			log.Printf("Failed to process row %s: %v", row.ID, err)
			continue
		}
		processed++
	}

	for _, row := range researchType {
		if err := s.processResearchType(row.Doc); err != nil {
			failed++
			log.Printf("Failed to process row %s: %v", row.ID, err)
			continue
		}
		processed++
	}

	for _, row := range manufacturers {
		if err := s.processManufacturer(row.Doc); err != nil {
			failed++
			log.Printf("Failed to process row %s: %v", row.ID, err)
			continue
		}
		processed++
	}

	for _, row := range clients {
		if err := s.processClient(row.Doc); err != nil {
			failed++
			log.Printf("Failed to process row %s: %v", row.ID, err)
			continue
		}
		processed++
	}

	for _, row := range contacts {
		if err := s.processContact(row.Doc); err != nil {
			failed++
			log.Printf("Failed to process row %s: %v", row.ID, err)
			continue
		}
		processed++
	}

	for _, row := range classificators {
		if err := s.processClassificator(row.Doc); err != nil {
			failed++
			log.Printf("Failed to process row %s: %v", row.ID, err)
			continue
		}
		processed++
	}

	for _, row := range devices {
		if err := s.processDevice(row.Doc); err != nil {
			failed++
			log.Printf("Failed to process row %s: %v", row.ID, err)
			continue
		}
		processed++
	}

	for _, row := range tickets {
		if err := s.processTicket(row.Doc); err != nil {
			failed++
			log.Printf("Failed to process row %s: %v", row.ID, err)
			continue
		}
		processed++
	}

	log.Printf("\nImport completed: %d processed, %d failed", processed, failed)
}

// Stream the CouchDB export file
func (s *ImportService) ImportFromFile(filePath string) error {
	log.Printf("Starting import from: %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode the entire export structure
	var export couchDBExport
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&export); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	log.Printf("Found %d total rows, processing %d rows",
		export.TotalRows, len(export.Rows))

	s.processRowsInOrder(export.Rows)

	return nil
}

func (s *ImportService) processDepartment(docBytes []byte) error {
	var proxy struct {
		models.Department
		RawID string `json:"_id"`
	}

	if err := json.Unmarshal(docBytes, &proxy); err != nil {
		return err
	}

	parsedID, err := extractAndParseUUID(proxy.RawID)

	if err != nil {
		return err
	}

	proxy.Department.ID = parsedID

	return s.departmentService.AddNewDepartment(context.Background(), proxy.Department)
}

func (s *ImportService) processUser(docBytes []byte, index int) error {
	var proxy struct {
		models.User
		Roles []string `json:"roles"`
		RawID string   `json:"_id"`
	}

	if err := json.Unmarshal(docBytes, &proxy); err != nil {
		return err
	}

	parsedID, err := extractAndParseUUID(proxy.RawID)

	if err != nil {
		return err
	}

	proxy.User.UserID = parsedID

	if len(proxy.Roles) == 0 {
		proxy.Roles = []string{"user"}
	} else if slices.Contains(proxy.Roles, "leader") {
		proxy.Roles = []string{"admin"}
	} else if slices.Contains(proxy.Roles, "coordinator") {
		proxy.Roles = []string{"coordinator"}
	} else {
		proxy.Roles = []string{"user"}
	}

	account, err := s.accountService.CreateUser(context.Background(), fmt.Sprintf("%s_%d", "user", index), "test123", proxy.Roles[0], &parsedID)

	if err != nil {
		return err
	}

	proxy.User.UserID = account.UserID

	return s.userService.CreateNewUser(context.Background(), proxy.User)
}

func (s *ImportService) processRegion(docBytes []byte) error {
	var proxy struct {
		repository.Region
		RawID string `json:"_id"`
	}

	if err := json.Unmarshal(docBytes, &proxy); err != nil {
		return err
	}

	parsedID, err := extractAndParseUUID(proxy.RawID)

	if err != nil {
		return err
	}

	proxy.Region.ID = parsedID

	return s.regionsRepo.AddNewRegion(context.Background(), proxy.Region)
}

func (s *ImportService) processResearchType(docBytes []byte) error {
	var proxy struct {
		repository.ResearchType
		RawID string `json:"_id"`
	}

	if err := json.Unmarshal(docBytes, &proxy); err != nil {
		return err
	}

	parsedID, err := extractAndParseUUID(proxy.RawID)

	if err != nil {
		return err
	}

	proxy.ResearchType.ID = parsedID

	return s.researchTypeRepo.AddNewResearchType(context.Background(), proxy.ResearchType)
}

func (s *ImportService) processManufacturer(docBytes []byte) error {
	var proxy struct {
		repository.Manufacturer
		RawID string `json:"_id"`
	}

	if err := json.Unmarshal(docBytes, &proxy); err != nil {
		return err
	}

	parsedID, err := extractAndParseUUID(proxy.RawID)

	if err != nil {
		return err
	}

	proxy.Manufacturer.ID = parsedID

	return s.manufacturerRepo.AddNewManufacturer(context.Background(), proxy.Manufacturer)
}

func (s *ImportService) processClient(docBytes []byte) error {
	var proxy struct {
		models.Client
		RawID string `json:"_id"`
	}

	if err := json.Unmarshal(docBytes, &proxy); err != nil {
		return err
	}

	parsedID, err := extractAndParseUUID(proxy.RawID)

	if err != nil {
		return err
	}

	proxy.Client.ID = parsedID

	return s.clientService.CreateClient(context.Background(), proxy.Client)
}

func (s *ImportService) processContact(docBytes []byte) error {
	var proxy struct {
		models.Contact
		Ref      uuid.UUID `json:"ref"`
		FullName string    `json:"firstName"`
		RawID    string    `json:"_id"`
	}

	if err := json.Unmarshal(docBytes, &proxy); err != nil {
		return err
	}

	parsedID, err := extractAndParseUUID(proxy.RawID)

	if err != nil {
		return err
	}

	proxy.Contact.ID = parsedID
	proxy.Contact.ClientID = proxy.Ref
	proxy.Contact.Name = proxy.FullName

	return s.contactService.CreateContact(context.Background(), proxy.Contact)
}

func (s *ImportService) processClassificator(docBytes []byte) error {
	type ManufacturerDetails struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	var proxy struct {
		models.Classificator
		Manufacturer *models.FlexibleManufacturer `json:"manufacturer"`
		RawID        string                       `json:"_id"`
	}

	if err := json.Unmarshal(docBytes, &proxy); err != nil {
		return err
	}

	parsedID, err := extractAndParseUUID(proxy.RawID)

	if err != nil {
		return err
	}

	proxy.Classificator.ID = parsedID
	if proxy.Manufacturer != nil {
		proxy.Classificator.Manufacturer = &proxy.Manufacturer.ID
	}

	return s.classificatorService.NewClassificator(context.Background(), proxy.Classificator)
}

func (s *ImportService) processDevice(docBytes []byte) error {
	var proxy struct {
		models.Device
		RawID string `json:"_id"`
	}

	if err := json.Unmarshal(docBytes, &proxy); err != nil {
		return err
	}

	parsedID, err := extractAndParseUUID(proxy.RawID)

	if err != nil {
		return err
	}

	proxy.Device.ID = parsedID

	return s.deviceService.CreateNewDevice(context.Background(), proxy.Device)
}

func (s *ImportService) processTicket(docBytes []byte) error {
	var proxy struct {
		models.RawTicket
		AssignedAt       *models.FlexibleTime `json:"assignedAt"`
		CreatedAt        *models.FlexibleTime `json:"createdAt"`
		PlannedInterval  models.Interval      `json:"plannedInterval"`
		AssignedInterval models.Interval      `json:"assignedInterval"`
		ActualInterval   models.Interval      `json:"actualInterval"`
		RawID            string               `json:"_id"`
	}

	if err := json.Unmarshal(docBytes, &proxy); err != nil {
		return err
	}

	parsedID, err := extractAndParseUUID(proxy.RawID)

	if err != nil {
		return err
	}

	proxy.RawTicket.ID = parsedID

	if proxy.CreatedAt != nil {
		proxy.RawTicket.CreatedAt = &proxy.CreatedAt.Time
	}

	if proxy.AssignedAt != nil {
		proxy.RawTicket.AssignedAt = &proxy.AssignedAt.Time
	}

	if proxy.PlannedInterval.Start != nil {
		proxy.RawTicket.PlannedStart = &proxy.PlannedInterval.Start.Time
	}

	if proxy.PlannedInterval.End != nil {
		proxy.RawTicket.PlannedEnd = &proxy.PlannedInterval.End.Time
	}

	if proxy.AssignedInterval.Start != nil {
		proxy.RawTicket.AssignedStart = &proxy.AssignedInterval.Start.Time
	}

	if proxy.AssignedInterval.End != nil {
		proxy.RawTicket.AssignedEnd = &proxy.AssignedInterval.End.Time
	}

	if proxy.ActualInterval.Start != nil {
		proxy.RawTicket.WorkStartedAt = &proxy.ActualInterval.Start.Time
	}

	if proxy.ActualInterval.End != nil {
		proxy.RawTicket.WorkFinishedAt = &proxy.ActualInterval.End.Time
		proxy.RawTicket.ClosedAt = &proxy.ActualInterval.End.Time
	}

	if proxy.RawTicket.Reason == "repairs" {
		proxy.RawTicket.Reason = "repair"
	}

	return s.ticketRepo.CreateRawTicket(proxy.RawTicket)
}
