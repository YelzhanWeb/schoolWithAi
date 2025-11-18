package handlers

import (
	"log"
	"net/http"
	"strconv"

	"backend/internal/domain/services"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseService *services.CourseService
}

func NewCourseHandler(courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.courseService.GetAllCourses(c.Request.Context())
	if err != nil {
		// ИСПРАВЛЕНО: логируем детали, клиенту отдаём общее сообщение
		log.Printf("Error fetching courses: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, modules, err := h.courseService.GetCourseDetails(c.Request.Context(), courseID)
	if err != nil {
		log.Printf("Error fetching course %d: %v", courseID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"course":  course,
		"modules": modules,
	})
}

func (h *CourseHandler) GetModuleResources(c *gin.Context) {
	moduleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	resources, err := h.courseService.GetModuleResources(c.Request.Context(), moduleID)
	if err != nil {
		log.Printf("Error fetching resources for module %d: %v", moduleID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resources": resources})
}
