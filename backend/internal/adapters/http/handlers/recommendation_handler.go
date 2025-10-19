package handlers

import (
	"backend/internal/domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	recommendationService *services.RecommendationService
}

func NewRecommendationHandler(recommendationService *services.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationService: recommendationService,
	}
}

func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recommendations, err := h.recommendationService.GetRecommendations(
		c.Request.Context(),
		userID.(int64),
		10,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
	})
}

func (h *RecommendationHandler) RefreshRecommendations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recommendations, err := h.recommendationService.RefreshRecommendations(
		c.Request.Context(),
		userID.(int64),
		10,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Recommendations refreshed",
		"recommendations": recommendations,
	})
}
