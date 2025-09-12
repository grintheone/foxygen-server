package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
)

type AttachmentService struct {
	repo repository.AttachmentRepository
}

func NewAttachmentService(repo repository.AttachmentRepository) *AttachmentService {
	return &AttachmentService{repo}
}

func (s *AttachmentService) generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)

	baseName := uuid.NewString()

	return baseName + ext
}

func (s *AttachmentService) UploadMultipleFiles(ctx context.Context, fileHeaders []*multipart.FileHeader, uploadDir string, refID uuid.UUID) ([]*models.Attachment, error) {
	var (
		attachments []*models.Attachment
		mu          sync.Mutex
		wg          sync.WaitGroup
		errors      = make(chan error, len(fileHeaders))
	)

	for _, fileHeader := range fileHeaders {
		wg.Add(1)

		go func(fh *multipart.FileHeader) {
			defer wg.Done()

			attachment, err := s.UploadFile(ctx, fh, uploadDir, refID)
			if err != nil {
				errors <- fmt.Errorf("failed to upload %s: %w", fh.Filename, err)
				return
			}

			mu.Lock()
			attachments = append(attachments, attachment)
			mu.Unlock()
		}(fileHeader)
	}

	wg.Wait()
	close(errors)

	// Collect all errors
	var uploadErrors []error
	for err := range errors {
		uploadErrors = append(uploadErrors, err)
	}

	if len(uploadErrors) > 0 {
		// If any upload failed, you might want to clean up successful ones
		if len(attachments) > 0 {
			// for _, attachment := range attachments {
			// 	s.DeleteFile(ctx, attachment.ID, uploadDir) // Clean up successful uploads
			// }
		}
		return nil, fmt.Errorf("multiple upload errors: %v", uploadErrors)
	}

	return attachments, nil
}

func (s *AttachmentService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, uploadDir string, refID uuid.UUID) (*models.Attachment, error) {
	// Generate unique file name
	originalName := fileHeader.Filename
	safeFileName := s.generateUniqueFileName(originalName)
	filePath := path.Join(uploadDir, safeFileName)

	// Open uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Get file info for size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Create attachment record
	attachment := &models.Attachment{
		FileName:     safeFileName,
		OriginalName: originalName,
		FilePath:     filePath,
		FileSize:     fileInfo.Size(),
		MimeType:     fileHeader.Header.Get("Content-Type"),
		RefID:        refID,
	}

	if err := s.repo.Create(ctx, attachment); err != nil {
		// Clean up file if database operation fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save attachment to database: %w", err)
	}

	return attachment, nil
}

func (s *AttachmentService) GetAttachmentsByRefID(ctx context.Context, refID uuid.UUID) (*[]models.Attachment, error) {
	attachments, err := s.repo.GetAttachmentsByRefID(ctx, refID)
	if err != nil {
		return nil, fmt.Errorf("service error fetching attachments: %w", err)
	}

	return attachments, nil
}

func (s *AttachmentService) GetAttachmentByID(ctx context.Context, id int) (*models.Attachment, error) {
	attachment, err := s.repo.GetAttachmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service error getting attachment by id: %w", err)
	}

	return attachment, nil
}
