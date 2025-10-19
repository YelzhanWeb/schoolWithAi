package handlers

import (
	"backend/internal/domain/models"
	"backend/internal/domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	progressService *services.ProgressService
}

func NewProgressHandler(progressService *services.ProgressService) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
	}
}

type UpdateProgressRequest struct {
	ResourceID int64  `json:"resource_id" binding:"required"`
	Status     string `json:"status" binding:"required"`
	Score      int    `json:"score"`
	TimeSpent  int    `json:"time_spent"`
	Rating     int    `json:"rating"` // опционально
}

// UpdateProgress - обновить прогресс урока
func (h *ProgressHandler) UpdateProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := &models.ProgressUpdate{
		StudentID:  userID.(int64),
		ResourceID: req.ResourceID,
		Status:     req.Status,
		Score:      req.Score,
		TimeSpent:  req.TimeSpent,
		Rating:     req.Rating,
	}

	if err := h.progressService.UpdateProgress(c.Request.Context(), update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Progress updated successfully",
	})
}

// GetMyProgress - получить свой прогресс
func (h *ProgressHandler) GetMyProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	progress, err := h.progressService.GetStudentProgress(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"progress": progress})
}

// GetMyStatistics - получить статистику
func (h *ProgressHandler) GetMyStatistics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	stats, err := h.progressService.GetStudentStatistics(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
