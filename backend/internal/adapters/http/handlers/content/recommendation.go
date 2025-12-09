package content

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *CourseHandler) GetRecommendations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		log.Error().Msg("unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courses, err := h.courseService.GetRecommendations(c.Request.Context(), userID.(string))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID.(string)).Msg("failed to get recommendations")
		return
	}

	c.JSON(http.StatusOK, courses)
}
