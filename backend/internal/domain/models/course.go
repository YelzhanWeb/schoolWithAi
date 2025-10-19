package models

import (
	"errors"
	"time"
)

type Course struct {
	ID              int64
	Title           string
	Description     string
	CreatedBy       int64
	DifficultyLevel int
	AgeGroup        string
	Subject         string
	IsPublished     bool
	ThumbnailURL    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Module struct {
	ID          int64
	CourseID    int64
	Title       string
	Description string
	OrderIndex  int
	CreatedAt   time.Time
}

type Resource struct {
	ID            int64
	ModuleID      int64
	Title         string
	Content       string
	ResourceType  string
	Difficulty    int
	EstimatedTime int
	FileURL       string
	ThumbnailURL  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Валидация на уровне домена
func (c *Course) Validate() error {
	if c.Title == "" {
		return errors.New("course title is required")
	}
	if c.DifficultyLevel < 1 || c.DifficultyLevel > 5 {
		return errors.New("difficulty level must be between 1 and 5")
	}
	return nil
}
