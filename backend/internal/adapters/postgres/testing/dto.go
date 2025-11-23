package testing

import (
	"time"

	"backend/internal/entities"
)

type testDTO struct {
	ID           string `db:"id"`
	ModuleID     string `db:"module_id"`
	Title        string `db:"title"`
	PassingScore int    `db:"passing_score"`
}

func (d *testDTO) toEntity() *entities.Test {
	return &entities.Test{
		ID:           d.ID,
		ModuleID:     d.ModuleID,
		Title:        d.Title,
		PassingScore: d.PassingScore,
		Questions:    []entities.Question{}, // Инициализируем пустым слайсом
	}
}

type questionDTO struct {
	ID           string `db:"id"`
	TestID       string `db:"test_id"`
	Text         string `db:"text"`
	QuestionType string `db:"question_type"`
}

func (d *questionDTO) toEntity() entities.Question {
	return entities.Question{
		ID:           d.ID,
		TestID:       d.TestID,
		Text:         d.Text,
		QuestionType: d.QuestionType,
		Answers:      []entities.Answer{},
	}
}

type answerDTO struct {
	ID         string `db:"id"`
	QuestionID string `db:"question_id"`
	Text       string `db:"text"`
	IsCorrect  bool   `db:"is_correct"`
}

func (d *answerDTO) toEntity() entities.Answer {
	return entities.Answer{
		ID:         d.ID,
		QuestionID: d.QuestionID,
		Text:       d.Text,
		IsCorrect:  d.IsCorrect,
	}
}

type resultDTO struct {
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"`
	TestID      string    `db:"test_id"`
	Score       int       `db:"score"`
	IsPassed    bool      `db:"is_passed"`
	AttemptDate time.Time `db:"attempt_date"`
}
