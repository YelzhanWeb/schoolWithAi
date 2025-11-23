package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID              string
	AuthorID        string
	SubjectID       string
	Title           string
	Description     string
	DifficultyLevel int
	Tags            []string
	CoverImageURL   string
	IsPublished     bool
	CreatedAt       time.Time

	Modules []Module
}

type Module struct {
	ID         string
	CourseID   string
	Title      string
	OrderIndex int

	// Опционально: список уроков
	Lessons []Lesson
}

type Lesson struct {
	ID                string
	ModuleID          string
	Title             string
	ContentText       string
	VideoURL          string
	FileAttachmentURL string
	XPReward          int
	OrderIndex        int
}

// Конструкторы

func NewCourse(authorID, subjectID, title string, difficulty int) (*Course, error) {
	if difficulty < 1 || difficulty > 5 {
		return nil, errors.New("difficulty must be between 1 and 5")
	}
	return &Course{
		ID:              uuid.NewString(),
		AuthorID:        authorID,
		SubjectID:       subjectID,
		Title:           title,
		DifficultyLevel: difficulty,
		Tags:            []string{},
		IsPublished:     false,
		CreatedAt:       time.Now().UTC(),
	}, nil
}

func NewModule(courseID, title string, order int) *Module {
	return &Module{
		ID:         uuid.NewString(),
		CourseID:   courseID,
		Title:      title,
		OrderIndex: order,
	}
}

func NewLesson(moduleID, title string, order int) *Lesson {
	return &Lesson{
		ID:         uuid.NewString(),
		ModuleID:   moduleID,
		Title:      title,
		XPReward:   10, // Default reward
		OrderIndex: order,
	}
}
