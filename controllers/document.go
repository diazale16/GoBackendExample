package controllers

import (
	"net/http"

	"github.com/example/supabase-migration-demo/internal/storage"
	"github.com/example/supabase-migration-demo/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentController struct {
	service *services.DocumentService
}

func NewDocumentController(service *services.DocumentService) *DocumentController {
	return &DocumentController{service: service}
}

func (c *DocumentController) Upload(ctx *gin.Context) {
	userIDStr := ctx.PostForm("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	allowedTypes := map[string]bool{
		"text/plain":               true,
		"image/png":                true,
		"image/jpeg":               true,
		"application/pdf":          true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] && header.Size > 10*1024*1024 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed or size exceeds 10MB"})
		return
	}

	doc, err := c.service.UploadDocument(userID, file, header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url := storage.BuildURL(doc.StorageKey)
	ctx.JSON(http.StatusCreated, doc.ToResponse(url))
}

func (c *DocumentController) GetAll(ctx *gin.Context) {
	docs, err := c.service.GetAllDocuments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]interface{}, len(docs))
	for i, doc := range docs {
		url := storage.BuildURL(doc.StorageKey)
		response[i] = doc.ToResponse(url)
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DocumentController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	doc, err := c.service.GetDocumentByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if doc == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	url := storage.BuildURL(doc.StorageKey)
	ctx.JSON(http.StatusOK, doc.ToResponse(url))
}

func (c *DocumentController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	if err := c.service.DeleteDocument(id); err != nil {
		if err.Error() == "document not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "document deleted successfully"})
}