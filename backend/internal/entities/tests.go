package entities

import (
	"time"

	"github.com/google/uuid"
)

type Test struct {
	ID           string
	ModuleID     string
	Title        string
	PassingScore int

	Questions []Question
}

type Question struct {
	ID           string
	TestID       string
	Text         string
	QuestionType string // 'single_choice', 'multiple'

	Answers []Answer
}

type Answer struct {
	ID         string
	QuestionID string
	Text       string
	IsCorrect  bool
}

type TestResult struct {
	ID          string
	UserID      string
	TestID      string
	Score       int
	IsPassed    bool
	AttemptDate time.Time
}

func NewTest(moduleID, title string, passingScore int) *Test {
	return &Test{
		ID:           uuid.NewString(),
		ModuleID:     moduleID,
		Title:        title,
		PassingScore: passingScore,
		Questions:    []Question{},
	}
}

func NewQuestion(testID, text, qType string) *Question {
	return &Question{
		ID:           uuid.NewString(),
		TestID:       testID,
		Text:         text,
		QuestionType: qType,
		Answers:      []Answer{},
	}
}

func NewAnswer(questionID, text string, isCorrect bool) *Answer {
	return &Answer{
		ID:         uuid.NewString(),
		QuestionID: questionID,
		Text:       text,
		IsCorrect:  isCorrect,
	}
}
