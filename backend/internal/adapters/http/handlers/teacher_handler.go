package handlers

import (
	"log"
	"net/http"
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/domain/services"

	"github.com/gin-gonic/gin"
)

type TeacherHandler struct {
	courseService  *services.CourseService
	teacherService *services.TeacherService
}

func NewTeacherHandler(courseService *services.CourseService, teacherService *services.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		courseService:  courseService,
		teacherService: teacherService,
	}
}

// ========================================
// УПРАВЛЕНИЕ КУРСАМИ
// ========================================

// GetMyCourses получает все курсы учителя
func (h *TeacherHandler) GetMyCourses(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courses, err := h.courseService.GetCoursesByTeacher(c.Request.Context(), teacherID)
	if err != nil {
		log.Printf("Error fetching courses for teacher %d: %v", teacherID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

type CreateCourseRequest struct {
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	DifficultyLevel int    `json:"difficulty_level" binding:"required,min=1,max=5"`
	AgeGroup        string `json:"age_group" binding:"required"`
	Subject         string `json:"subject" binding:"required"`
}

// CreateCourse создает новый курс
func (h *TeacherHandler) CreateCourse(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course := &models.Course{
		Title:           req.Title,
		Description:     req.Description,
		CreatedBy:       teacherID,
		DifficultyLevel: req.DifficultyLevel,
		AgeGroup:        req.AgeGroup,
		Subject:         req.Subject,
		IsPublished:     false,
		ThumbnailURL:    "/images/default.jpg",
	}

	if err := h.courseService.CreateCourse(c.Request.Context(), course); err != nil {
		log.Printf("Error creating course: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Course created successfully",
		"course":  course,
	})
}

// UpdateCourse обновляет курс
func (h *TeacherHandler) UpdateCourse(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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
		CreatedBy:       teacherID,
		DifficultyLevel: req.DifficultyLevel,
		AgeGroup:        req.AgeGroup,
		Subject:         req.Subject,
	}

	if err := h.courseService.UpdateCourse(c.Request.Context(), course); err != nil {
		log.Printf("Error updating course: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

// DeleteCourse удаляет курс
func (h *TeacherHandler) DeleteCourse(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if err := h.teacherService.DeleteCourse(c.Request.Context(), courseID, teacherID); err != nil {
		log.Printf("Error deleting course: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

// PublishCourse публикует курс
func (h *TeacherHandler) PublishCourse(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if err := h.teacherService.PublishCourse(c.Request.Context(), courseID, teacherID); err != nil {
		log.Printf("Error publishing course: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course published successfully"})
}

// UnpublishCourse снимает курс с публикации
func (h *TeacherHandler) UnpublishCourse(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if err := h.teacherService.UnpublishCourse(c.Request.Context(), courseID, teacherID); err != nil {
		log.Printf("Error unpublishing course: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course unpublished successfully"})
}

// ========================================
// УПРАВЛЕНИЕ МОДУЛЯМИ
// ========================================

// CreateModule создает модуль
func (h *TeacherHandler) CreateModule(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	module := &models.Module{
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
	}

	if err := h.teacherService.CreateModule(c.Request.Context(), module, teacherID); err != nil {
		log.Printf("Error creating module: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Module created successfully",
		"module":  module,
	})
}

// UpdateModule обновляет модуль
func (h *TeacherHandler) UpdateModule(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	moduleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	var req models.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	module := &models.Module{
		ID:          moduleID,
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
	}

	if err := h.teacherService.UpdateModule(c.Request.Context(), module, teacherID); err != nil {
		log.Printf("Error updating module: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module updated successfully"})
}

// DeleteModule удаляет модуль
func (h *TeacherHandler) DeleteModule(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	moduleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	courseID, err := strconv.ParseInt(c.Query("course_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID is required"})
		return
	}

	if err := h.teacherService.DeleteModule(c.Request.Context(), moduleID, courseID, teacherID); err != nil {
		log.Printf("Error deleting module: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module deleted successfully"})
}

// ========================================
// УПРАВЛЕНИЕ РЕСУРСАМИ
// ========================================

// CreateResource создает ресурс
func (h *TeacherHandler) CreateResource(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resource := &models.Resource{
		ModuleID:      req.ModuleID,
		Title:         req.Title,
		Content:       req.Content,
		ResourceType:  req.ResourceType,
		Difficulty:    req.Difficulty,
		EstimatedTime: req.EstimatedTime,
		FileURL:       req.FileURL,
		ThumbnailURL:  req.ThumbnailURL,
	}

	if err := h.teacherService.CreateResource(c.Request.Context(), resource, teacherID); err != nil {
		log.Printf("Error creating resource: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Resource created successfully",
		"resource": resource,
	})
}

// UpdateResource обновляет ресурс
func (h *TeacherHandler) UpdateResource(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	resourceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	var req models.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	moduleID, err := strconv.ParseInt(c.Query("module_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Module ID is required"})
		return
	}

	resource := &models.Resource{
		ID:            resourceID,
		ModuleID:      moduleID,
		Title:         req.Title,
		Content:       req.Content,
		ResourceType:  req.ResourceType,
		Difficulty:    req.Difficulty,
		EstimatedTime: req.EstimatedTime,
		FileURL:       req.FileURL,
		ThumbnailURL:  req.ThumbnailURL,
	}

	if err := h.teacherService.UpdateResource(c.Request.Context(), resource, teacherID); err != nil {
		log.Printf("Error updating resource: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource updated successfully"})
}

// DeleteResource удаляет ресурс
func (h *TeacherHandler) DeleteResource(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	resourceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	moduleID, err := strconv.ParseInt(c.Query("module_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Module ID is required"})
		return
	}

	if err := h.teacherService.DeleteResource(c.Request.Context(), resourceID, moduleID, teacherID); err != nil {
		log.Printf("Error deleting resource: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}

// ========================================
// СТАТИСТИКА
// ========================================

// GetDashboard получает общую статистику учителя
func (h *TeacherHandler) GetDashboard(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	stats, err := h.teacherService.GetDashboardStatistics(c.Request.Context(), teacherID)
	if err != nil {
		log.Printf("Error fetching dashboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetCourseStatistics получает статистику по курсу
func (h *TeacherHandler) GetCourseStatistics(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	stats, err := h.teacherService.GetCourseStatistics(c.Request.Context(), courseID, teacherID)
	if err != nil {
		log.Printf("Error fetching course statistics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetCourseStudents получает список студентов и их прогресс по курсу
func (h *TeacherHandler) GetCourseStudents(c *gin.Context) {
	teacherID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	students, err := h.teacherService.GetCourseStudentProgress(c.Request.Context(), courseID, teacherID)
	if err != nil {
		log.Printf("Error fetching student progress: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"students": students})
}

// getUserID безопасно извлекает user_id
func getUserID(c *gin.Context) (int64, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		return 0, false
	}

	return userID, true
}
