package handlers

import (
	"context"
	"net/http"

	"backend/internal/entities"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

type SubjectResponse struct {
	ID     string `json:"id"`
	Slug   string `json:"slug"`
	NameRU string `json:"name_ru"`
	NameKZ string `json:"name_kz"`
}

type GetAllSubjectsResponse struct {
	Subjects []SubjectResponse `json:"subjects"`
}

// GetAllSubjects godoc
// @Summary Get all subjects
// @Description Get list of available subjects
// @Tags subjects
// @Produce json
// @Success 200 {array} GetAllSubjectsResponse
// @Failure 500
// @Router /v1/subjects [get]
func (h *SubjectHandler) GetAllSubjects(c *gin.Context) {
	subjects, err := h.service.GetAllSubjects(c.Request.Context())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Msg("failed to get all subjects")
		return
	}

	subs := make([]SubjectResponse, 0, len(subjects))
	for _, s := range subjects {
		subs = append(subs, SubjectResponse{
			ID:     s.ID,
			Slug:   s.Slug,
			NameRU: s.NameRu,
			NameKZ: s.NameKz,
		})
	}

	c.JSON(http.StatusOK, GetAllSubjectsResponse{
		Subjects: subs,
	})

	log.Info().Msg("all subjects got successfully")
}
