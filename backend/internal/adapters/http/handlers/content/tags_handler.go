package content

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type GetAllTagsResponse struct {
	Tags []TagResponse `json:"tags"`
}

// GetTags godoc
// @Summary Get all tags
// @Description Get list of available tags
// @Tags tags
// @Produce json
// @Success 200 {object} GetAllTagsResponse
// @Failure 500
// @Router /v1/tags [get]
func (h *CourseHandler) GetTags(c *gin.Context) {
	tags, err := h.courseService.GetAllTags(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get all tags")
		c.Status(http.StatusInternalServerError)
		return
	}

	tagsResp := make([]TagResponse, 0, len(tags))
	for _, t := range tags {
		tagsResp = append(tagsResp, TagResponse{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		})
	}

	c.JSON(http.StatusOK, GetAllTagsResponse{
		Tags: tagsResp,
	})

	log.Info().Msg("all tags got successfully")
}
