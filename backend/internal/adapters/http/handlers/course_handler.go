package handlers

import (
	"context"
	"net/http"

	"backend/internal/entities"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CourseService interface {
	CreateCourse(ctx context.Context, course *entities.Course) error
	GetUserCourse(ctx context.Context, courseID, userID string) (*entities.Course, error)
	UpdateCourse(ctx context.Context, courseID string, updates *entities.Course) error
	ChangePublishStatus(ctx context.Context, courseID string, isPublished bool) error

	CreateModule(ctx context.Context, userID string, module *entities.Module) error
	UpdateModule(ctx context.Context, userID string, module *entities.Module) error
	DeleteModule(ctx context.Context, userID, moduleID string) error

	CreateLesson(ctx context.Context, userID string, lesson *entities.Lesson) error
	UpdateLesson(ctx context.Context, userID string, lesson *entities.Lesson) error
	DeleteLesson(ctx context.Context, userID, lessonID string) error

	GetFullStructure(ctx context.Context, courseID string) ([]entities.Module, error)
}

type CourseHandler struct {
	courseService CourseService
}

func NewCourseHandler(service CourseService) *CourseHandler {
	return &CourseHandler{courseService: service}
}

type CreateCourseRequest struct {
	Title           string `json:"title"            binding:"required"`
	Description     string `json:"description"`
	SubjectID       string `json:"subject_id"       binding:"required"`
	DifficultyLevel int    `json:"difficulty_level" binding:"required,min=1,max=5"`
}

type CreateCourseResponse struct {
	ID string `json:"id"`
}

// CreateCourse godoc
// @Summary Create a new course
// @Description Create a course (Teacher/Admin only)
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body CreateCourseRequest true "Course data"
// @Success 201 {object} CreateCourseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500
// @Router /v1/courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		log.Error().Msg("user unauthorized")
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "unauthorized"})
		return
	}

	if role != "teacher" && role != "admin" {
		log.Error().Msg("only teachers can create courses")
		c.JSON(http.StatusForbidden, ErrorResponse{Message: "only teachers can create courses"})
		return
	}

	userID, _ := c.Get("user_id")

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse json request")
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	course, err := entities.NewCourse(userID.(string), req.SubjectID, req.Title, req.DifficultyLevel)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID.(string)).Msg("failed to create new course entity")
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	course.Description = req.Description

	err = h.courseService.CreateCourse(c.Request.Context(), course)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("course_id", course.ID).Msg("failed to create course")
		return
	}

	c.JSON(http.StatusCreated, CreateCourseResponse{ID: course.ID})
	log.Info().
		Str("user_id", userID.(string)).
		Str("course_id", course.ID).
		Str("title", course.Title).
		Int("difficulty_level", course.DifficultyLevel).
		Msg("course created successfully")
}

type UpdateCourseRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	DifficultyLevel int    `json:"difficulty_level"`
	CoverImageURL   string `json:"cover_image_url"`
}

// UpdateCourse godoc
// @Summary Update course details
// @Description Update title, description, cover, etc. (Author only)
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param input body UpdateCourseRequest true "Update data"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500
// @Router /v1/courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	courseID := c.Param("id")
	userID := c.GetString("user_id")

	_, err := h.courseService.GetUserCourse(c.Request.Context(), courseID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{Message: err.Error()})
		return
	}

	var req UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid input"})
		log.Error().Err(err).Str("user_id", userID).Str("course_id", courseID).Msg("failed to parse json request")
		return
	}

	updates := &entities.Course{
		Title:           req.Title,
		Description:     req.Description,
		DifficultyLevel: req.DifficultyLevel,
		CoverImageURL:   req.CoverImageURL,
	}

	if err := h.courseService.UpdateCourse(c.Request.Context(), courseID, updates); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID).Str("course_id", courseID).Msg("failed to update course")
		return
	}

	c.Status(http.StatusOK)
	log.Info().
		Str("user_id", userID).
		Str("course_id", courseID).
		Str("title", updates.Title).
		Int("difficulty_level", updates.DifficultyLevel).
		Msg("course updated successfully")
}

type PublishStatusRequest struct {
	IsPublished bool `json:"is_published"`
}

// ChangePublishStatus godoc
// @Summary Publish or Unpublish course
// @Description Change visibility of the course
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param input body PublishStatusRequest true "Status"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500
// @Router /v1/courses/{id}/publish [post]
func (h *CourseHandler) ChangePublishStatus(c *gin.Context) {
	courseID := c.Param("id")
	userID := c.GetString("user_id")

	_, err := h.courseService.GetUserCourse(c.Request.Context(), courseID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{Message: err.Error()})
		return
	}

	var req PublishStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid input"})
		log.Error().Err(err).Str("user_id", userID).Str("course_id", courseID).Msg("failed to parse json request")
		return
	}

	if err := h.courseService.ChangePublishStatus(c.Request.Context(), courseID, req.IsPublished); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID).Str("course_id", courseID).Msg("failed to published course")
		return
	}

	c.Status(http.StatusOK)
	log.Info().
		Str("user_id", userID).
		Str("course_id", courseID).
		Msg("course published successfully")
}

type CreateModuleRequest struct {
	CourseID   string `json:"course_id"   binding:"required"`
	Title      string `json:"title"       binding:"required"`
	OrderIndex int    `json:"order_index" binding:"required"`
}

type CreateModuleResponse struct {
	ModuleID string `json:"module_id" binding:"required"`
}

// CreateModule godoc
// @Summary Add module to course
// @Tags modules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body CreateModuleRequest true "Module data"
// @Success 201 {object} CreateModuleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500
// @Router /v1/modules [post]
func (h *CourseHandler) CreateModule(c *gin.Context) {
	userID := c.GetString("user_id")
	var req CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		log.Error().Err(err).Str("user_id", userID).Msg("failed to parse json request")
		return
	}

	module := entities.NewModule(req.CourseID, req.Title, req.OrderIndex)

	if err := h.courseService.CreateModule(c.Request.Context(), userID, module); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID).Msg("failed to create module")
		return
	}

	c.JSON(http.StatusCreated, CreateModuleResponse{
		ModuleID: module.ID,
	})
	log.Info().
		Str("user_id", userID).
		Str("course_id", module.CourseID).
		Str("module_id", module.ID).
		Str("module_title", module.Title).
		Int("module_order_index", module.OrderIndex).
		Msg("module created successfully")
}

type UpdateModuleRequest struct {
	Title      string `json:"title"`
	OrderIndex int    `json:"order_index"`
}

// UpdateModule godoc
// @Summary Update module
// @Tags modules
// @Security BearerAuth
// @Accept json
// @Param id path string true "Module ID"
// @Param input body UpdateModuleRequest true "Data"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500
// @Router /v1/modules/{id} [put]
func (h *CourseHandler) UpdateModule(c *gin.Context) {
	userID := c.GetString("user_id")
	moduleID := c.Param("id")

	var req UpdateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		log.Error().Err(err).Str("user_id", userID).Str("module_id", moduleID).Msg("failed to parse json request")
		return
	}

	module := &entities.Module{
		ID:         moduleID,
		Title:      req.Title,
		OrderIndex: req.OrderIndex,
	}

	if err := h.courseService.UpdateModule(c.Request.Context(), userID, module); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID).Str("module_id", moduleID).Msg("failed to update module")
		return
	}

	c.Status(http.StatusOK)
	log.Info().
		Str("user_id", userID).
		Str("course_id", module.CourseID).
		Str("module_id", module.ID).
		Str("module_title", module.Title).
		Int("module_order_index", module.OrderIndex).
		Msg("module updated successfully")
}

// DeleteModule godoc
// @Summary Delete module
// @Tags modules
// @Security BearerAuth
// @Param id path string true "Module ID"
// @Success 200
// @Failure 500
// @Router /v1/modules/{id} [delete]
func (h *CourseHandler) DeleteModule(c *gin.Context) {
	userID := c.GetString("user_id")
	moduleID := c.Param("id")

	if err := h.courseService.DeleteModule(c.Request.Context(), userID, moduleID); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID).Str("module_id", moduleID).Msg("failed to delete module")
		return
	}

	c.Status(http.StatusOK)
	log.Info().
		Str("user_id", userID).
		Str("module_id", moduleID).
		Msg("module deleted successfully")
}

type CreateLessonRequest struct {
	ModuleID    string `json:"module_id"    binding:"required"`
	Title       string `json:"title"        binding:"required"`
	ContentText string `json:"content_text"`
	VideoURL    string `json:"video_url"`
	OrderIndex  int    `json:"order_index"  binding:"required"`
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

// GetStructure godoc
// @Summary Get full course structure
// @Description Get modules and lessons for editor
// @Tags courses
// @Security BearerAuth
// @Param id path string true "Course ID"
// @Success 200 {object} GetStructureResponse
// @Failure 500
// @Router /v1/courses/{id}/structure [get]
func (h *CourseHandler) GetStructure(c *gin.Context) {
	courseID := c.Param("id")

	modulesEntities, err := h.courseService.GetFullStructure(c.Request.Context(), courseID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("course_id", courseID).Msg("failed to get course structure")
		return
	}

	modulesResp := make([]ModuleResponse, 0, len(modulesEntities))

	for _, m := range modulesEntities {
		lessonsResp := make([]LessonResponse, 0, len(m.Lessons))
		for _, l := range m.Lessons {
			lessonsResp = append(lessonsResp, LessonResponse{
				ID:                l.ID,
				ModuleID:          l.ModuleID,
				Title:             l.Title,
				ContentText:       l.ContentText,
				VideoURL:          l.VideoURL,
				FileAttachmentURL: l.FileAttachmentURL,
				XPReward:          l.XPReward,
				OrderIndex:        l.OrderIndex,
			})
		}

		modulesResp = append(modulesResp, ModuleResponse{
			ID:         m.ID,
			CourseID:   m.CourseID,
			Title:      m.Title,
			OrderIndex: m.OrderIndex,
			Lessons:    lessonsResp,
		})
	}

	c.JSON(http.StatusOK, GetStructureResponse{
		Modules: modulesResp,
	})
	log.Info().
		Str("course_id", courseID).
		Msg("course structure got successfully")
}
