package content

import (
	"context"
	"errors"
	"net/http"

	"backend/internal/entities"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CourseService interface {
	CreateCourse(ctx context.Context, course *entities.Course) error
	GetUserCourse(ctx context.Context, courseID, userID string) (*entities.Course, error)
	GetCourseByID(ctx context.Context, courseID string) (*entities.Course, error)
	UpdateCourse(ctx context.Context, courseID string, updates *entities.Course) error
	ChangePublishStatus(ctx context.Context, courseID string, isPublished bool) error
	GetCoursesByAuthor(ctx context.Context, authorID string) ([]entities.Course, error)
	DeleteCourse(ctx context.Context, id string) error
	GetCatalog(ctx context.Context) ([]entities.Course, error)

	CreateModule(ctx context.Context, userID string, module *entities.Module) error
	UpdateModule(ctx context.Context, userID string, module *entities.Module) error
	DeleteModule(ctx context.Context, userID, moduleID string) error

	CreateLesson(ctx context.Context, userID string, lesson *entities.Lesson) error
	GetLessonByID(ctx context.Context, lessonID string) (*entities.Lesson, error)
	UpdateLesson(ctx context.Context, userID string, lesson *entities.Lesson) error
	DeleteLesson(ctx context.Context, userID, lessonID string) error

	GetFullStructure(ctx context.Context, courseID string) ([]entities.Module, error)

	GetAllTags(ctx context.Context) ([]entities.Tag, error)
}

type ErrorResponse struct {
	Message string `json:"message" example:"something went wrong"`
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
	CoverImageURL   string `json:"cover_image_url"`
	Tags            []int  `json:"tags"`
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
	course.CoverImageURL = req.CoverImageURL

	for _, tagID := range req.Tags {
		course.Tags = append(course.Tags, entities.Tag{ID: tagID})
	}

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

// GetCatalog godoc
// @Summary Get all published courses
// @Tags courses
// @Produce json
// @Success 200 {object} CourseListResponse
// @Router /v1/catalog [get]
func (h *CourseHandler) GetCatalog(c *gin.Context) {
	courses, err := h.courseService.GetCatalog(c.Request.Context())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to get catalog")
		return
	}

	respCourses := make([]CourseDetailResponse, 0, len(courses))
	for _, course := range courses {
		tagsResp := make([]TagResponse, 0, len(course.Tags))
		for _, t := range course.Tags {
			tagsResp = append(tagsResp, TagResponse{
				ID:   t.ID,
				Name: t.Name,
				Slug: t.Slug,
			})
		}
		// ----------------------------------------

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

type TagResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CourseDetailResponse struct {
	ID              string        `json:"id"`
	AuthorID        string        `json:"author_id"`
	SubjectID       string        `json:"subject_id"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	DifficultyLevel int           `json:"difficulty_level"`
	CoverImageURL   string        `json:"cover_image_url"`
	IsPublished     bool          `json:"is_published"`
	Tags            []TagResponse `json:"tags"`
}

// GetCourse godoc
// @Summary Get course details
// @Description Get details of a specific course by ID
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} CourseDetailResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/courses/{id} [get]
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseID := c.Param("id")

	course, err := h.courseService.GetCourseByID(c.Request.Context(), courseID)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "course not found"})
			return
		}
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("course_id", courseID).Msg("failed to get course")
		return
	}

	tagsResp := make([]TagResponse, 0, len(course.Tags))
	for _, t := range course.Tags {
		tagsResp = append(tagsResp, TagResponse{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		})
	}

	resp := CourseDetailResponse{
		ID:              course.ID,
		AuthorID:        course.AuthorID,
		SubjectID:       course.SubjectID,
		Title:           course.Title,
		Description:     course.Description,
		DifficultyLevel: course.DifficultyLevel,
		CoverImageURL:   course.CoverImageURL,
		IsPublished:     course.IsPublished,
		Tags:            tagsResp,
	}

	c.JSON(http.StatusOK, resp)
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

type CourseListResponse struct {
	Courses []CourseDetailResponse `json:"courses"`
}

// GetMyCourses godoc
// @Summary Get courses created by current teacher
// @Tags courses
// @Security BearerAuth
// @Produce json
// @Success 200 {object} CourseListResponse
// @Failure 500
// @Router /teacher/courses [get]
func (h *CourseHandler) GetMyCourses(c *gin.Context) {
	userID := c.GetString("user_id") // Берем из JWT токена

	courses, err := h.courseService.GetCoursesByAuthor(c.Request.Context(), userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID).Msg("failed to get my courses")
		return
	}

	respCourses := make([]CourseDetailResponse, 0, len(courses))
	for _, course := range courses {
		respCourses = append(respCourses, CourseDetailResponse{
			ID:              course.ID,
			AuthorID:        course.AuthorID,
			SubjectID:       course.SubjectID,
			Title:           course.Title,
			Description:     course.Description,
			DifficultyLevel: course.DifficultyLevel,
			CoverImageURL:   course.CoverImageURL,
			IsPublished:     course.IsPublished,
		})
	}

	c.JSON(http.StatusOK, CourseListResponse{Courses: respCourses})
	log.Info().Str("author_id", userID).Msg("author courses got successfully")
}

type UpdateCourseRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	DifficultyLevel int    `json:"difficulty_level"`
	CoverImageURL   string `json:"cover_image_url"`
	SubjectID       string `json:"subject_id"`
	Tags            []int  `json:"tags"`
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
		SubjectID:       req.SubjectID,
	}

	for _, tagID := range req.Tags {
		updates.Tags = append(updates.Tags, entities.Tag{ID: tagID})
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

// DeleteCourse godoc
// @Summary Delete course
// @Description Delete course by id
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500
// @Router /v1/courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	courseID := c.Param("id")
	userID := c.GetString("user_id")

	_, err := h.courseService.GetUserCourse(c.Request.Context(), courseID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{Message: err.Error()})
		return
	}

	if err := h.courseService.DeleteCourse(c.Request.Context(), courseID); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("user_id", userID).Str("course_id", courseID).Msg("failed to delete course")
		return
	}

	c.Status(http.StatusOK)
	log.Info().
		Str("user_id", userID).
		Str("course_id", courseID).
		Msg("course deleted successfully")
}
