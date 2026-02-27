package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/grintheone/foxygen-server/internal/models"
	"github.com/grintheone/foxygen-server/internal/repository"
	"github.com/minio/minio-go/v7"
)

type AttachmentService struct {
	repo       repository.AttachmentRepository
	storage    *minio.Client
	bucketName string
}

func NewAttachmentService(repo repository.AttachmentRepository, storage *minio.Client, bucketName string) *AttachmentService {
	return &AttachmentService{
		repo:       repo,
		storage:    storage,
		bucketName: bucketName,
	}
}

func (s *AttachmentService) generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)

	baseName := uuid.NewString()

	return baseName + ext
}

func (s *AttachmentService) UploadMultipleFiles(ctx context.Context, fileHeaders []*multipart.FileHeader, refID uuid.UUID) ([]*models.Attachment, error) {
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

			attachment, err := s.UploadFile(ctx, fh, refID)
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

func (s *AttachmentService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, refID uuid.UUID) (*models.Attachment, error) {
	// Generate object name
	originalName := fileHeader.Filename
	objectName := s.generateUniqueFileName(originalName)

	// Open uploaded file stream
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Upload object to MinIO
	if _, err := s.storage.PutObject(
		ctx,
		s.bucketName,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")},
	); err != nil {
		return nil, fmt.Errorf("failed to upload object to minio: %w", err)
	}

	// Create attachment record
	ext := filepath.Ext(originalName)
	attachment := &models.Attachment{
		ID:        objectName,
		Name:      originalName,
		MediaType: fileHeader.Header.Get("Content-Type"),
		Ext:       ext,
		RefID:     refID,
	}

	if err := s.repo.Create(ctx, attachment); err != nil {
		_ = s.storage.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
		return nil, fmt.Errorf("failed to save attachment to database: %w", err)
	}

	return attachment, nil
}

func (s *AttachmentService) GetFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	object, err := s.storage.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from minio: %w", err)
	}

	if _, err := object.Stat(); err != nil {
		_ = object.Close()
		return nil, fmt.Errorf("failed to stat object from minio: %w", err)
	}

	return object, nil
}

func (s *AttachmentService) GetAttachmentsByRefID(ctx context.Context, refID uuid.UUID) ([]*models.Attachment, error) {
	attachments, err := s.repo.GetAttachmentsByRefID(ctx, refID)
	if err != nil {
		return nil, fmt.Errorf("service error fetching attachments: %w", err)
	}

	return attachments, nil
}

func (s *AttachmentService) GetAttachmentByID(ctx context.Context, id string) (*models.Attachment, error) {
	attachment, err := s.repo.GetAttachmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service error getting attachment by id: %w", err)
	}

	return attachment, nil
}
