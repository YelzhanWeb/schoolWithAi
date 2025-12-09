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

	respCourses := make([]CourseDetailResponse, 0, len(courses))
	for _, course := range courses {
		// tags
		tagsResp := make([]TagResponse, 0, len(course.Tags))
		for _, t := range course.Tags {
			tagsResp = append(tagsResp, TagResponse{
				ID:   t.ID,
				Name: t.Name,
				Slug: t.Slug,
			})
		}

		// author (если есть)
		var authorResp *AuthorResponse
		if course.Author != nil {
			authorResp = &AuthorResponse{
				ID:        course.Author.ID,
				FullName:  course.Author.FirstName,
				AvatarURL: course.Author.AvatarURL,
			}
		}

		// cover image: если в entities.Course CoverImageURL типа *string — разкомментируй соответствующую обработку ниже.
		coverImage := course.CoverImageURL
		/*
			// Если в Course используется *string:
			if course.CoverImageURL != nil {
				coverImage = *course.CoverImageURL
			} else {
				coverImage = ""
			}
		*/

		respCourses = append(respCourses, CourseDetailResponse{
			ID:              course.ID,
			AuthorID:        course.AuthorID,
			SubjectID:       course.SubjectID,
			Title:           course.Title,
			Description:     course.Description,
			DifficultyLevel: course.DifficultyLevel,
			CoverImageURL:   coverImage,
			IsPublished:     course.IsPublished,
			Tags:            tagsResp,

			Author: authorResp,
		})
	}

	c.JSON(http.StatusOK, CourseListResponse{Courses: respCourses})
}
