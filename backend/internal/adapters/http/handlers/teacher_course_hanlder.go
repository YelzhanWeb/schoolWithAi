package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/domain/services"

	"github.com/gin-gonic/gin"
)

type TeacherCourseHandler struct {
	courseService *services.CourseService
	// Мы будем использовать тот же CourseService,
	// но добавим в него методы для учителя
}

func NewTeacherCourseHandler(courseService *services.CourseService) *TeacherCourseHandler {
	return &TeacherCourseHandler{
		courseService: courseService,
	}
}

// GetMyCourses получает курсы, созданные текущим учителем
func (h *TeacherCourseHandler) GetMyCourses(c *gin.Context) {
	userID, _ := c.Get("user_id") // Получаем из AuthMiddleware

	courses, err := h.courseService.GetCoursesByTeacher(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

// Структура для создания/обновления курса
type CreateCourseRequest struct {
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	DifficultyLevel int    `json:"difficulty_level" binding:"required,min=1,max=5"`
	AgeGroup        string `json:"age_group" binding:"required"`
	Subject         string `json:"subject" binding:"required"`
}

// CreateCourse создает новый курс
func (h *TeacherCourseHandler) CreateCourse(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := &models.Course{
		Title:           req.Title,
		Description:     req.Description,
		CreatedBy:       userID.(int64), // Устанавливаем создателя
		DifficultyLevel: req.DifficultyLevel,
		AgeGroup:        req.AgeGroup,
		Subject:         req.Subject,
		IsPublished:     false,                 // По умолчанию не опубликован
		ThumbnailURL:    "/images/default.jpg", // Заглушка
	}

	if err := h.courseService.CreateCourse(c.Request.Context(), course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// UpdateCourse обновляет существующий курс
func (h *TeacherCourseHandler) UpdateCourse(c *gin.Context) {
	userID, _ := c.Get("user_id")
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := &models.Course{
		ID:              courseID,
		Title:           req.Title,
		Description:     req.Description,
		CreatedBy:       userID.(int64), // Для проверки прав
		DifficultyLevel: req.DifficultyLevel,
		AgeGroup:        req.AgeGroup,
		Subject:         req.Subject,
		// IsPublished не меняем здесь, сделаем отдельный эндпоинт
	}

	if err := h.courseService.UpdateCourse(c.Request.Context(), course); err != nil {
		if err.Error() == "course not found or permission denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course updated"})
}

// TODO: Добавить хендлеры для CreateModule, CreateResource, PublishCourse...
