package handlers

import (
	"context"
	"net/http"

	"backend/internal/entities"

	"github.com/gin-gonic/gin"
)

type CourseService interface {
	CreateCourse(ctx context.Context, course *entities.Course) error
	GetUserCourse(ctx context.Context, courseID, userID string) (*entities.Course, error)
	UpdateCourse(ctx context.Context, courseID string, updates *entities.Course) error
	ChangePublishStatus(ctx context.Context, courseID string, isPublished bool) error
}

type CourseHandler struct {
	courseService CourseService
}

func NewCourseHandler(service CourseService) *CourseHandler {
	return &CourseHandler{courseService: service}
}

type CreateCourseRequest struct {
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	SubjectID       string `json:"subject_id" binding:"required"`
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
// @Failure 500 {object} ErrorResponse
// @Router /v1/courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "unauthorized"})
		return
	}

	if role != "teacher" && role != "admin" {
		c.JSON(http.StatusForbidden, ErrorResponse{Message: "only teachers can create courses"})
		return
	}

	userID, _ := c.Get("user_id")

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	course, err := entities.NewCourse(userID.(string), req.SubjectID, req.Title, req.DifficultyLevel)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	course.Description = req.Description

	err = h.courseService.CreateCourse(c.Request.Context(), course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateCourseResponse{ID: course.ID})
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
// @Failure 500 {object} ErrorResponse
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
		return
	}

	updates := &entities.Course{
		Title:           req.Title,
		Description:     req.Description,
		DifficultyLevel: req.DifficultyLevel,
		CoverImageURL:   req.CoverImageURL,
	}

	if err := h.courseService.UpdateCourse(c.Request.Context(), courseID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusOK)
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
// @Failure 500 {object} ErrorResponse
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
		return
	}

	if err := h.courseService.ChangePublishStatus(c.Request.Context(), courseID, req.IsPublished); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
