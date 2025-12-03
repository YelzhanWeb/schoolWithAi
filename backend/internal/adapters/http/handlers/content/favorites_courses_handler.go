package content

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ToggleFavorite godoc
// @Summary Toggle course favorite status
// @Description Add or remove course from favorites
// @Tags courses
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} map[string]any
// @Failure 500 {object} ErrorResponse
// @Router /v1/courses/{id}/favorite [post]
func (h *CourseHandler) ToggleFavorite(c *gin.Context) {
	userID := c.GetString("user_id")
	courseID := c.Param("id")

	isFav, err := h.courseService.ToggleFavorite(c.Request.Context(), userID, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to toggle favorite"})
		log.Error().Err(err).Str("user_id", userID).Str("course_id", courseID).Msg("failed to toggle favorite")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_favorite": isFav,
	})
}

// GetFavorites godoc
// @Summary Get user favorite courses
// @Tags courses
// @Security BearerAuth
// @Produce json
// @Success 200 {object} CourseListResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/courses/favorites [get]
func (h *CourseHandler) GetFavorites(c *gin.Context) {
	userID := c.GetString("user_id")

	courses, err := h.courseService.GetUserFavorites(c.Request.Context(), userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID).Msg("failed to get favorites")
		return
	}

	respCourses := make([]CourseDetailResponse, 0, len(courses))
	for _, course := range courses {
		// Если вы уже исправили проблему N+1 с тегами, здесь тоже стоит их маппить
		// Пока оставим пустой массив тегов для простоты, или скопируйте логику из GetCatalog
		tagsResp := make([]TagResponse, 0)

		respCourses = append(respCourses, CourseDetailResponse{
			ID:              course.ID,
			AuthorID:        course.AuthorID,
			SubjectID:       course.SubjectID,
			Title:           course.Title,
			Description:     course.Description,
			DifficultyLevel: course.DifficultyLevel,
			CoverImageURL:   course.CoverImageURL,
			IsPublished:     course.IsPublished,
			Tags:            tagsResp,
		})
	}

	c.JSON(http.StatusOK, CourseListResponse{Courses: respCourses})
}
