// backend/internal/adapters/http/handlers/leaderboard_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/services/student"

	"github.com/gin-gonic/gin"
)

type LeaderboardHandler struct {
	service *student.StudentService
}

func NewLeaderboardHandler(service *student.StudentService) *LeaderboardHandler {
	return &LeaderboardHandler{service: service}
}

type LeaderboardEntry struct {
	Rank      int    `json:"rank"`
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
	XP        int64  `json:"xp"`
	Level     int    `json:"level"`
	LeagueID  int    `json:"league_id,omitempty"`
}

type LeaderboardResponse struct {
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	UserRank    *int               `json:"user_rank,omitempty"`
}

// GetWeeklyLeaderboard godoc
// @Summary Get weekly league leaderboard
// @Description Get top players in current user's league for this week
// @Tags leaderboard
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Number of entries (default 50)"
// @Success 200 {object} LeaderboardResponse
// @Router /v1/leaderboard/weekly [get]
func (h *LeaderboardHandler) GetWeeklyLeaderboard(c *gin.Context) {
	userID := c.GetString("user_id")

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	leaderboard, userRank, err := h.service.GetWeeklyLeaderboard(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get leaderboard"})
		return
	}

	c.JSON(http.StatusOK, LeaderboardResponse{
		Leaderboard: leaderboard,
		UserRank:    userRank,
	})
}

// GetGlobalLeaderboard godoc
// @Summary Get global leaderboard
// @Description Get top players by total XP (all-time)
// @Tags leaderboard
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Number of entries (default 50)"
// @Success 200 {object} LeaderboardResponse
// @Router /v1/leaderboard/global [get]
func (h *LeaderboardHandler) GetGlobalLeaderboard(c *gin.Context) {
	userID := c.GetString("user_id")

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	leaderboard, userRank, err := h.service.GetGlobalLeaderboard(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get leaderboard"})
		return
	}

	c.JSON(http.StatusOK, LeaderboardResponse{
		Leaderboard: leaderboard,
		UserRank:    userRank,
	})
}
