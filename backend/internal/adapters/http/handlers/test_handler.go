package handlers

import (
	"context"
	"errors"
	"net/http"

	"backend/internal/entities"
	"backend/internal/services/student"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type TestService interface {
	CreateFullTest(ctx context.Context, test *entities.Test) error
	GetTestByModule(ctx context.Context, moduleID string) (*entities.Test, error)
	UpdateFullTest(ctx context.Context, test *entities.Test) error
	DeleteTest(ctx context.Context, testID string) error
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
	ID        string `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type QuestionResponse struct {
	ID           string           `json:"id"`
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
// @Success 200 {object} TestResponse
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
	moduleID := c.Param("id")
	test, err := h.service.GetTestByModule(c.Request.Context(), moduleID)
	if err != nil {
		log.Error().Err(err).Str("module_id", moduleID).Msg("failed to get test by moduleID")
		if errors.Is(err, entities.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, mapTestToResponse(test, false))
	log.Info().Str("module_id", moduleID).Str("test_id", test.ID).Msg("test got successfully")
}

// GetTestWithAnswer godoc
// @Summary Get test details with correct answer
// @Description Get test structure with questions and answers
// @Tags tests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Success 200 {object} TestResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404
// @Failure 500
// @Router /v1/modules/{id}/test-with-answers [get]
func (h *TestHandler) GetTestWithAnswer(c *gin.Context) {
	_, exists := c.Get("role")
	if !exists {
		log.Error().Msg("user unauthorized")
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "unauthorized"})
		return
	}

	moduleID := c.Param("id")
	test, err := h.service.GetTestByModule(c.Request.Context(), moduleID)
	if err != nil {
		log.Error().Err(err).Str("module_id", moduleID).Msg("failed to get test by moduleID")
		if errors.Is(err, entities.ErrNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, mapTestToResponse(test, true))
	log.Info().Str("module_id", moduleID).Str("test_id", test.ID).Msg("test got successfully")
}

func mapTestToResponse(test *entities.Test, withCorrectAnswer bool) TestResponse {
	resp := TestResponse{
		TestID:       test.ID,
		ModuleID:     test.ModuleID,
		Title:        test.Title,
		PassingScore: test.PassingScore,
	}

	for _, q := range test.Questions {
		qResp := QuestionResponse{
			ID:           q.ID,
			Text:         q.Text,
			QuestionType: q.QuestionType,
		}

		for _, a := range q.Answers {
			if withCorrectAnswer {
				qResp.Answers = append(qResp.Answers, AnswerResponse{
					ID:        a.ID,
					Text:      a.Text,
					IsCorrect: a.IsCorrect,
				})
				continue
			}
			qResp.Answers = append(qResp.Answers, AnswerResponse{
				ID:   a.ID,
				Text: a.Text,
			})
		}

		resp.Questions = append(resp.Questions, qResp)
	}

	return resp
}

// UpdateTest godoc
// @Summary Update old test
// @Description Update a test (Teacher/Admin only)
// @Tags tests
// @Security BearerAuth
// @Accept json
// @Param input body CreateTestRequest true "Test data"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500
// @Router /v1/tests/:id [put]
func (h *TestHandler) UpdateTest(c *gin.Context) {
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
	testID := c.Param("id")
	var req CreateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	test := entities.NewTest(req.ModuleID, req.Title, req.PassingScore)
	test.ID = testID

	for _, qReq := range req.Questions {
		question := entities.NewQuestion(test.ID, qReq.Text, qReq.QuestionType)
		for _, aReq := range qReq.Answers {
			answer := entities.NewAnswer(question.ID, aReq.Text, aReq.IsCorrect)
			question.Answers = append(question.Answers, *answer)
		}
		test.Questions = append(test.Questions, *question)
	}

	if err := h.service.UpdateFullTest(c.Request.Context(), test); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("test_id", testID).Msg("failed to update test")
		return
	}

	c.Status(http.StatusOK)
	log.Info().Str("test_id", testID).Msg("test updated successfully")
}

// DeleteTest godoc
// @Summary Delete old test
// @Description Delete a test (Teacher/Admin only)
// @Tags tests
// @Security BearerAuth
// @Success 200
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500
// @Router /v1/tests/:id [delete]
func (h *TestHandler) DeleteTest(c *gin.Context) {
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
	testID := c.Param("id")

	if err := h.service.DeleteTest(c.Request.Context(), testID); err != nil {
		c.Status(http.StatusInternalServerError)
		log.Error().Err(err).Str("test_id", testID).Msg("failed to delete test")
		return
	}
	c.Status(http.StatusOK)
	log.Info().Str("test_id", testID).Msg("test deleted successfully")
}

type SubmitTestRequest struct {
	TestID  string `json:"test_id" binding:"required"`
	Answers []struct {
		QuestionID string `json:"question_id"`
		AnswerID   string `json:"answer_id"`
	} `json:"answers"`
}

func (h *StudentHandler) SubmitTest(c *gin.Context) {
	userID := c.GetString("user_id")
	var req SubmitTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	// Мапим в структуру сервиса
	srvAnswers := make([]student.StudentAnswer, len(req.Answers))
	for i, a := range req.Answers {
		srvAnswers[i] = student.StudentAnswer{QuestionID: a.QuestionID, AnswerID: a.AnswerID}
	}

	res, xp, err := h.service.SubmitTest(c.Request.Context(), userID, req.TestID, srvAnswers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_passed": res.IsPassed,
		"score":     res.Score,
		"xp_gained": xp,
	})
}
