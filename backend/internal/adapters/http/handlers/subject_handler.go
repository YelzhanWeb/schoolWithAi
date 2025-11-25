package handlers

import (
	"context"
	"net/http"

	"backend/internal/entities"

	"github.com/gin-gonic/gin"
)

type SubjectService interface {
	GetAllSubjects(ctx context.Context) ([]entities.Subject, error)
}

type SubjectHandler struct {
	service SubjectService
}

func NewSubjectHandler(service SubjectService) *SubjectHandler {
	return &SubjectHandler{service: service}
}

// GetAllSubjects godoc
// @Summary Get all subjects
// @Description Get list of available subjects
// @Tags subjects
// @Produce json
// @Success 200 {array} entities.Subject
// @Failure 500 {object} ErrorResponse
// @Router /v1/subjects [get]
func (h *SubjectHandler) GetAllSubjects(c *gin.Context) {
	subjects, err := h.service.GetAllSubjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subjects": subjects})
}
