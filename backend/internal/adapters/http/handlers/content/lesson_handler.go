package content

import (
	"errors"
	"net/http"

	"backend/internal/entities"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreateLessonRequest struct {
	ModuleID          string `json:"module_id"           binding:"required"`
	Title             string `json:"title"               binding:"required"`
	ContentText       string `json:"content_text"`
	VideoURL          string `json:"video_url"`
	FileAttachmentURL string `json:"file_attachment_url"`
	OrderIndex        int    `json:"order_index"         binding:"required"`
	XPReward          int    `json:"xp_reward"`
}

type CreateLessonResponse struct {
	LessonID string `json:"lesson_id" binding:"required"`
}

// CreateLesson godoc
// @Summary Add lesson to module
// @Tags lessons
// @Security BearerAuth
// @Accept json
// @Param input body CreateLessonRequest true "Lesson data"
// @Success 201 {object} CreateLessonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500
// @Router /v1/lessons [post]
func (h *CourseHandler) CreateLesson(c *gin.Context) {
	userID := c.GetString("user_id")
	var req CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	lesson := entities.NewLesson(req.ModuleID, req.Title, req.OrderIndex)
	lesson.ContentText = req.ContentText
	lesson.VideoURL = req.VideoURL
	lesson.FileAttachmentURL = req.FileAttachmentURL

	if req.XPReward > 0 {
		lesson.XPReward = req.XPReward
	}

	if err := h.courseService.CreateLesson(c.Request.Context(), userID, lesson); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID).Str("module_id", req.ModuleID).Msg("failed to create lesson")
		return
	}

	c.JSON(http.StatusCreated, CreateLessonResponse{
		LessonID: lesson.ID,
	})

	log.Info().
		Str("user_id", userID).
		Str("module_id", req.ModuleID).
		Str("lesson_id", lesson.ID).
		Str("lesson_title", lesson.Title).
		Msg("module created successfully")
}

type LessonResponse struct {
	ID                string `json:"id"`
	ModuleID          string `json:"module_id"`
	Title             string `json:"title"`
	ContentText       string `json:"content_text"`
	VideoURL          string `json:"video_url"`
	FileAttachmentURL string `json:"file_attachment_url"`
	XPReward          int    `json:"xp_reward"`
	OrderIndex        int    `json:"order_index"`
}

// GetLesson godoc
// @Summary Get full lesson details
// @Description Get content, video url and attachments for a specific lesson
// @Tags lessons
// @Security BearerAuth
// @Param id path string true "Lesson ID"
// @Success 200 {object} LessonResponse
// @Failure 404
// @Failure 500
// @Router /v1/lessons/{id} [get]
func (h *CourseHandler) GetLesson(c *gin.Context) {
	lessonID := c.Param("id")

	lesson, err := h.courseService.GetLessonByID(c.Request.Context(), lessonID)
	if err != nil {
		log.Error().Err(err).Str("lesson_id", lessonID).Msg("failed to get lesson by id")
		if errors.Is(err, entities.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, LessonResponse{
		ID:                lesson.ID,
		ModuleID:          lesson.ModuleID,
		Title:             lesson.Title,
		ContentText:       lesson.ContentText,
		VideoURL:          lesson.VideoURL,
		FileAttachmentURL: lesson.FileAttachmentURL,
		XPReward:          lesson.XPReward,
		OrderIndex:        lesson.OrderIndex,
	})
}

type UpdateLessonRequest struct {
	Title             string `json:"title"`
	ContentText       string `json:"content_text"`
	VideoURL          string `json:"video_url"`
	FileAttachmentURL string `json:"file_attachment_url"`
	OrderIndex        int    `json:"order_index"`
	XPReward          int    `json:"xp_reward"`
}

// UpdateLesson godoc
// @Summary Update lesson content
// @Tags lessons
// @Security BearerAuth
// @Accept json
// @Param id path string true "Lesson ID"
// @Param input body UpdateLessonRequest true "Data"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500
// @Router /v1/lessons/{id} [put]
func (h *CourseHandler) UpdateLesson(c *gin.Context) {
	userID := c.GetString("user_id")
	lessonID := c.Param("id")

	var req UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid json request"})
		log.Error().Err(err).Str("user_id", userID).Str("lesson_id", lessonID).Msg("failed to parse json request")
		return
	}
	lesson := &entities.Lesson{
		ID:                lessonID,
		Title:             req.Title,
		ContentText:       req.ContentText,
		VideoURL:          req.VideoURL,
		FileAttachmentURL: req.FileAttachmentURL,
		OrderIndex:        req.OrderIndex,
	}

	if req.XPReward > 0 {
		lesson.XPReward = req.XPReward
	}

	if err := h.courseService.UpdateLesson(c.Request.Context(), userID, lesson); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID).Str("lesson_id", lessonID).Msg("failed to update lesson")
		return
	}
	c.Status(http.StatusOK)
	log.Info().Str("user_id", userID).Str("lesson_id", lessonID).Msg("lesson updated successfully")
}

// DeleteLesson godoc
// @Summary Delete lesson content
// @Tags lessons
// @Security BearerAuth
// @Param id path string true "Lesson ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500
// @Router /v1/lessons/{id} [delete]
func (h *CourseHandler) DeleteLesson(c *gin.Context) {
	userID := c.GetString("user_id")
	lessonID := c.Param("id")

	if err := h.courseService.DeleteLesson(c.Request.Context(), userID, lessonID); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID).Str("lesson_id", lessonID).Msg("failed to delete lesson")
		return
	}
	c.Status(http.StatusOK)
	log.Info().Str("user_id", userID).Str("lesson_id", lessonID).Msg("lesson deleted successfully")
}

type GetStructureResponse struct {
	Modules []ModuleResponse `json:"modules"`
}

type ModuleResponse struct {
	ID         string           `json:"id"`
	CourseID   string           `json:"course_id"`
	Title      string           `json:"title"`
	OrderIndex int              `json:"order_index"`
	Lessons    []LessonResponse `json:"lessons"`
}
