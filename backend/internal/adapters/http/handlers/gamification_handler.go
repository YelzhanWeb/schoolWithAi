package handlers

import (
	"context"
	"net/http"

	"backend/internal/entities"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type GamificationService interface {
	GetAllLeagues(ctx context.Context) ([]entities.League, error)
}

type GamificationHandler struct {
	service GamificationService
}

func NewGamificationHandler(service GamificationService) *GamificationHandler {
	return &GamificationHandler{service: service}
}

type LeagueResponse struct {
	ID         int    `json:"id"`
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	OrderIndex int    `json:"order_index"`
	IconURL    string `json:"icon_url"`
}

type LeaguesListResponse struct {
	Leagues []LeagueResponse `json:"leagues"`
}

// GetAllLeagues godoc
// @Summary Get all available leagues
// @Tags gamification
// @Produce json
// @Success 200 {object} LeaguesListResponse
// @Router /v1/gamification/leagues [get]
func (h *GamificationHandler) GetAllLeagues(c *gin.Context) {
	leagues, err := h.service.GetAllLeagues(c.Request.Context())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to get leagues")
		return
	}

	var response []LeagueResponse
	for _, l := range leagues {
		response = append(response, LeagueResponse{
			ID:         l.ID,
			Slug:       l.Slug,
			Name:       l.Name,
			OrderIndex: l.OrderIndex,
			IconURL:    l.IconURL,
		})
	}

	c.JSON(http.StatusOK, LeaguesListResponse{Leagues: response})
}
