package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/service"
)

type OpinionHandler struct {
	opinionService service.OpinionService
}

func NewOpinionHandler(opinionService service.OpinionService) *OpinionHandler {
	return &OpinionHandler{opinionService: opinionService}
}

type CreateOpinionRequest struct {
	WriterID      uint64  `json:"writer_id"                binding:"required"`
	WorkID        uint64  `json:"work_id"                  binding:"required"`
	Sentiment     bool    `json:"sentiment"                binding:"required"`
	Quote         string  `json:"quote"                    binding:"required"`
	Source        string  `json:"source"                   binding:"required"`
	Page          *string `json:"page,omitempty"`
	StatementYear *int    `json:"statement_year,omitempty"`
}

type UpdateOpinionRequest struct {
	Sentiment     bool    `json:"sentiment"                binding:"required"`
	Quote         string  `json:"quote"                    binding:"required"`
	Source        string  `json:"source"                   binding:"required"`
	Page          *string `json:"page,omitempty"`
	StatementYear *int    `json:"statement_year,omitempty"`
}

func (h *OpinionHandler) Create(c *gin.Context) {
	var req CreateOpinionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	opinion, err := h.opinionService.CreateOpinion(
		req.WriterID,
		req.WorkID,
		req.Sentiment,
		req.Quote,
		req.Source,
		req.Page,
		req.StatementYear,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"writer_id":      opinion.WriterID(),
		"work_id":        opinion.WorkID(),
		"sentiment":      opinion.Sentiment(),
		"quote":          opinion.Quote(),
		"source":         opinion.Source(),
		"page":           opinion.Page(),
		"statement_year": opinion.StatementYear(),
	})
}

func (h *OpinionHandler) GetByWriter(c *gin.Context) {
	writerIDStr := c.Param("writer_id")
	writerID, err := strconv.ParseUint(writerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid writer_id"})
		return
	}

	opinions, err := h.opinionService.GetOpinionsByWriter(writerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.opinionsToResponse(opinions))
}

func (h *OpinionHandler) GetByWork(c *gin.Context) {
	workIDStr := c.Param("work_id")
	workID, err := strconv.ParseUint(workIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_id"})
		return
	}

	opinions, err := h.opinionService.GetOpinionsByWork(workID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.opinionsToResponse(opinions))
}

func (h *OpinionHandler) GetByWriterAndWork(c *gin.Context) {
	writerIDStr := c.Param("writer_id")
	writerID, err := strconv.ParseUint(writerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid writer_id"})
		return
	}

	workIDStr := c.Param("work_id")
	workID, err := strconv.ParseUint(workIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_id"})
		return
	}

	opinion, err := h.opinionService.GetOpinion(writerID, workID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "opinion not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"writer_id":      opinion.WriterID(),
		"work_id":        opinion.WorkID(),
		"sentiment":      opinion.Sentiment(),
		"quote":          opinion.Quote(),
		"source":         opinion.Source(),
		"page":           opinion.Page(),
		"statement_year": opinion.StatementYear(),
	})
}

func (h *OpinionHandler) List(c *gin.Context) {
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

	opinions, err := h.opinionService.ListOpinions(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.opinionsToResponse(opinions))
}

func (h *OpinionHandler) opinionsToResponse(opinions []*domain.Opinion) []gin.H {
	result := make([]gin.H, len(opinions))
	for i, o := range opinions {
		result[i] = gin.H{
			"writer_id":      o.WriterID(),
			"work_id":        o.WorkID(),
			"sentiment":      o.Sentiment(),
			"quote":          o.Quote(),
			"source":         o.Source(),
			"page":           o.Page(),
			"statement_year": o.StatementYear(),
		}
	}
	return result
}

func (h *OpinionHandler) Update(c *gin.Context) {
	writerIDStr := c.Param("writer_id")
	writerID, err := strconv.ParseUint(writerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid writer_id"})
		return
	}

	workIDStr := c.Param("work_id")
	workID, err := strconv.ParseUint(workIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_id"})
		return
	}

	var req UpdateOpinionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.opinionService.UpdateOpinion(writerID, workID, req.Sentiment, req.Quote, req.Source, req.Page, req.StatementYear); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "opinion updated"})
}

func (h *OpinionHandler) Delete(c *gin.Context) {
	writerIDStr := c.Param("writer_id")
	writerID, err := strconv.ParseUint(writerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid writer_id"})
		return
	}

	workIDStr := c.Param("work_id")
	workID, err := strconv.ParseUint(workIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_id"})
		return
	}

	if err := h.opinionService.DeleteOpinion(writerID, workID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "opinion deleted"})
}
