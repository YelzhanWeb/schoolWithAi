package content

import (
	"net/http"

	"backend/internal/entities"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

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
