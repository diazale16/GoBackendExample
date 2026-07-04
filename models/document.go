package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	FileName   string         `gorm:"type:varchar(255);not null" json:"file_name"`
	StorageKey string         `gorm:"type:varchar(512);not null" json:"storage_key"`
	MimeType   string         `gorm:"type:varchar(100)" json:"mime_type"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

type DocumentResponse struct {
	ID         uuid.UUID `json:"id"`
	FileName   string    `json:"file_name"`
	StorageKey string    `json:"storage_key"`
	MimeType   string    `json:"mime_type"`
	URL        string    `json:"url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (d *Document) ToResponse(url string) DocumentResponse {
	return DocumentResponse{
		ID:         d.ID,
		FileName:   d.FileName,
		StorageKey: d.StorageKey,
		MimeType:   d.MimeType,
		URL:        url,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}

type UploadDocumentRequest struct {
	FileName string `form:"file_name"`
}