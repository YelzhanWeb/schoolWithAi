// backend/internal/adapters/http/handlers/resource_handler.go - СОЗДАТЬ!
package handlers

import (
	"backend/internal/domain/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResourceHandler struct {
	courseService *services.CourseService
}

func NewResourceHandler(courseService *services.CourseService) *ResourceHandler {
	return &ResourceHandler{courseService: courseService}
}

// GetResource - получить детали ресурса
func (h *ResourceHandler) GetResource(c *gin.Context) {
	resourceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	// НУЖНО ДОБАВИТЬ МЕТОД В CourseService
	resource, err := h.courseService.GetResourceByID(c.Request.Context(), resourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resource": resource})
}
