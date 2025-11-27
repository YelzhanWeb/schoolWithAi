package handlers

import (
	"context"
	"errors"
	"net/http"

	"backend/internal/entities"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type TestService interface {
	CreateFullTest(ctx context.Context, test *entities.Test) error
	GetTestByModule(ctx context.Context, moduleID string) (*entities.Test, error)
}

type TestHandler struct {
	service TestService
}

func NewTestHandler(service TestService) *TestHandler {
	return &TestHandler{service: service}
}

type CreateAnswerRequest struct {
	Text      string `json:"text" binding:"required"`
	IsCorrect bool   `json:"is_correct"`
}

type CreateQuestionRequest struct {
	Text         string                `json:"text" binding:"required"`
	QuestionType string                `json:"question_type" binding:"required"` // single_choice
	Answers      []CreateAnswerRequest `json:"answers" binding:"required,min=2"`
}

type CreateTestRequest struct {
	ModuleID     string                  `json:"module_id" binding:"required"`
	Title        string                  `json:"title" binding:"required"`
	PassingScore int                     `json:"passing_score" binding:"required"`
	Questions    []CreateQuestionRequest `json:"questions" binding:"required,min=1"`
}

type CreateTestResponse struct {
	TestID string `json:"test_id"`
}

// CreateTest godoc
// @Summary Create a new test
// @Description Create a test (Teacher/Admin only)
// @Tags tests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body CreateTestRequest true "Test data"
// @Success 201 {object} CreateTestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500
// @Router /v1/tests [post]
func (h *TestHandler) CreateTest(c *gin.Context) {
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

	var req CreateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse json request")
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	test := entities.NewTest(req.ModuleID, req.Title, req.PassingScore)

	for _, qReq := range req.Questions {
		question := entities.NewQuestion(test.ID, qReq.Text, qReq.QuestionType)

		for _, aReq := range qReq.Answers {
			answer := entities.NewAnswer(question.ID, aReq.Text, aReq.IsCorrect)
			question.Answers = append(question.Answers, *answer)
		}
		test.Questions = append(test.Questions, *question)
	}

	if err := h.service.CreateFullTest(c.Request.Context(), test); err != nil {
		log.Error().Err(err).Str("module_id", test.ModuleID).Str("test_id", test.ID).Msg("failed to create test")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, CreateTestResponse{
		TestID: test.ID,
	})

	log.Info().Str("module_id", test.ModuleID).Str("test_id", test.ID).Msg("test created successfully")
}

type AnswerResponse struct {
	Text string `json:"text"`
	// IsCorrect bool   `json:"is_correct"`
}

type QuestionResponse struct {
	Text         string           `json:"text"`
	QuestionType string           `json:"question_type"`
	Answers      []AnswerResponse `json:"answers"`
}

type TestResponse struct {
	TestID       string             `json:"test_id"`
	ModuleID     string             `json:"module_id"`
	Title        string             `json:"title"`
	PassingScore int                `json:"passing_score"`
	Questions    []QuestionResponse `json:"questions"`
}

// GetTest godoc
// @Summary Get test details
// @Description Get test structure with questions and answers
// @Tags tests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Success 201 {object} CreateTestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404
// @Failure 500
// @Router /v1/modules/{id}/test [get]
func (h *TestHandler) GetTest(c *gin.Context) {
	_, exists := c.Get("role")
	if !exists {
		log.Error().Msg("user unauthorized")
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "unauthorized"})
		return
	}
	moduleID := c.Param("moduleId")
	test, err := h.service.GetTestByModule(c.Request.Context(), moduleID)
	if err != nil {
		log.Error().Err(err).Str("module_id", moduleID).Msg("failed to get test by moduleID")
		if errors.Is(err, entities.ErrNotFound) {
			c.Status(http.StatusNotFound)
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, mapTestToResponse(test))
	log.Info().Str("module_id", moduleID).Str("test_id", test.ID).Msg("test got successfully")
}

func mapTestToResponse(test *entities.Test) TestResponse {
	resp := TestResponse{
		TestID:       test.ID,
		ModuleID:     test.ModuleID,
		Title:        test.Title,
		PassingScore: test.PassingScore,
	}

	for _, q := range test.Questions {
		qResp := QuestionResponse{
			Text:         q.Text,
			QuestionType: q.QuestionType,
		}

		for _, a := range q.Answers {
			qResp.Answers = append(qResp.Answers, AnswerResponse{
				Text: a.Text,
			})
		}

		resp.Questions = append(resp.Questions, qResp)
	}

	return resp
}
