package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/what-writers-like/backend/internal/service"
)

type WriterHandler struct {
	writerService service.WriterService
}

func NewWriterHandler(writerService service.WriterService) *WriterHandler {
	return &WriterHandler{writerService: writerService}
}

type CreateWriterRequest struct {
	Name      string  `json:"name"                 binding:"required"`
	BirthYear int     `json:"birth_year"           binding:"required"`
	DeathYear *int    `json:"death_year,omitempty"`
	Bio       *string `json:"bio,omitempty"`
}

type UpdateWriterRequest struct {
	Name      string  `json:"name"                 binding:"required"`
	BirthYear int     `json:"birth_year"           binding:"required"`
	DeathYear *int    `json:"death_year,omitempty"`
	Bio       *string `json:"bio,omitempty"`
}

func (h *WriterHandler) Create(c *gin.Context) {
	var req CreateWriterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	writer, err := h.writerService.CreateWriter(req.Name, req.BirthYear, req.DeathYear, req.Bio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         writer.ID(),
		"name":       writer.Name(),
		"birth_year": writer.BirthYear(),
		"death_year": writer.DeathYear(),
		"bio":        writer.Bio(),
	})
}

func (h *WriterHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	writer, err := h.writerService.GetWriter(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "writer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         writer.ID(),
		"name":       writer.Name(),
		"birth_year": writer.BirthYear(),
		"death_year": writer.DeathYear(),
		"bio":        writer.Bio(),
	})
}

func (h *WriterHandler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	writers, err := h.writerService.ListWriters(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]gin.H, len(writers))
	for i, w := range writers {
		result[i] = gin.H{
			"id":         w.ID(),
			"name":       w.Name(),
			"birth_year": w.BirthYear(),
			"death_year": w.DeathYear(),
			"bio":        w.Bio(),
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *WriterHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateWriterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.writerService.UpdateWriter(id, req.Name, req.BirthYear, req.DeathYear, req.Bio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "writer updated"})
}

func (h *WriterHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.writerService.DeleteWriter(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "writer deleted"})
}
