package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/what-writers-like/backend/internal/service"
)

type WorkHandler struct {
	workService service.WorkService
}

func NewWorkHandler(workService service.WorkService) *WorkHandler {
	return &WorkHandler{workService: workService}
}

type CreateWorkRequest struct {
	Title    string `json:"title"     binding:"required"`
	AuthorID uint64 `json:"author_id" binding:"required"`
}

type UpdateWorkRequest struct {
	Title    string `json:"title"     binding:"required"`
	AuthorID uint64 `json:"author_id" binding:"required"`
}

func (h *WorkHandler) Create(c *gin.Context) {
	var req CreateWorkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	work, err := h.workService.CreateWork(req.Title, req.AuthorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        work.ID(),
		"title":     work.Title(),
		"author_id": work.AuthorID(),
	})
}

func (h *WorkHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	work, err := h.workService.GetWork(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "work not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        work.ID(),
		"title":     work.Title(),
		"author_id": work.AuthorID(),
	})
}

func (h *WorkHandler) GetByAuthor(c *gin.Context) {
	authorIDStr := c.Param("author_id")
	authorID, err := strconv.ParseUint(authorIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author_id"})
		return
	}

	works, err := h.workService.GetWorksByAuthor(authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]gin.H, len(works))
	for i, w := range works {
		result[i] = gin.H{
			"id":        w.ID(),
			"title":     w.Title(),
			"author_id": w.AuthorID(),
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *WorkHandler) List(c *gin.Context) {
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

	works, err := h.workService.ListWorks(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]gin.H, len(works))
	for i, w := range works {
		result[i] = gin.H{
			"id":        w.ID(),
			"title":     w.Title(),
			"author_id": w.AuthorID(),
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *WorkHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateWorkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.workService.UpdateWork(id, req.Title, req.AuthorID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "work updated"})
}

func (h *WorkHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.workService.DeleteWork(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "work deleted"})
}
