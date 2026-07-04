package repositories

import (
	"fmt"

	"github.com/example/supabase-migration-demo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(doc *models.Document) error {
	result := r.db.Create(doc)
	if result.Error != nil {
		return fmt.Errorf("failed to create document: %w", result.Error)
	}
	return nil
}

func (r *DocumentRepository) GetAll() ([]models.Document, error) {
	var docs []models.Document
	result := r.db.Find(&docs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get documents: %w", result.Error)
	}
	return docs, nil
}

func (r *DocumentRepository) GetByID(id uuid.UUID) (*models.Document, error) {
	var doc models.Document
	result := r.db.First(&doc, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get document: %w", result.Error)
	}
	return &doc, nil
}

func (r *DocumentRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Document{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete document: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("document not found")
	}
	return nil
}

func (r *DocumentRepository) GetByUserID(userID uuid.UUID) ([]models.Document, error) {
	var docs []models.Document
	result := r.db.Where("user_id = ?", userID).Find(&docs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get documents by user: %w", result.Error)
	}
	return docs, nil
}