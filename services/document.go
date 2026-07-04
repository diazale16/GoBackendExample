package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/example/supabase-migration-demo/internal/storage"
	"github.com/example/supabase-migration-demo/models"
	"github.com/example/supabase-migration-demo/repositories"
	"github.com/google/uuid"
)

type DocumentService struct {
	docRepo *repositories.DocumentRepository
	storage *storage.S3Client
}

func NewDocumentService(
	docRepo *repositories.DocumentRepository,
	storage *storage.S3Client,
) *DocumentService {
	return &DocumentService{
		docRepo: docRepo,
		storage: storage,
	}
}

func (s *DocumentService) UploadDocument(
	file multipart.File,
	header *multipart.FileHeader,
) (*models.Document, error) {
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	uploadResult, err := s.storage.Upload(
		context.Background(),
		file,
		header.Filename,
		contentType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to storage: %w", err)
	}

	doc := &models.Document{
		FileName:   header.Filename,
		StorageKey: uploadResult.Key,
		MimeType:   uploadResult.MimeType,
	}

	if err := s.docRepo.Create(doc); err != nil {
		if delErr := s.storage.Delete(context.Background(), uploadResult.Key); delErr != nil {
			fmt.Printf("Warning: failed to cleanup storage after DB error: %v\n", delErr)
		}
		return nil, fmt.Errorf("failed to save document metadata: %w", err)
	}

	return doc, nil
}

func (s *DocumentService) GetAllDocuments() ([]models.Document, error) {
	return s.docRepo.GetAll()
}

func (s *DocumentService) GetDocumentByID(id uuid.UUID) (*models.Document, error) {
	return s.docRepo.GetByID(id)
}

func (s *DocumentService) DeleteDocument(id uuid.UUID) error {
	doc, err := s.docRepo.GetByID(id)
	if err != nil {
		return err
	}
	if doc == nil {
		return fmt.Errorf("document not found")
	}

	if err := s.docRepo.Delete(id); err != nil {
		return err
	}

	if err := s.storage.Delete(context.Background(), doc.StorageKey); err != nil {
		fmt.Printf("Warning: document deleted from DB but failed to delete from storage: %v\n", err)
	}

	return nil
}

func (s *DocumentService) GetFileContent(file io.Reader, header *multipart.FileHeader) (io.Reader, error) {
	return file, nil
}