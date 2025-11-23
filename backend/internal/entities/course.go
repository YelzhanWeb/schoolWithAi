package entities

import (
	"errors"
	"time"
)

type Course struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	CreatedBy       int64     `json:"created_by"`
	DifficultyLevel int       `json:"difficulty_level"`
	AgeGroup        string    `json:"age_group"`
	Subject         string    `json:"subject"`
	IsPublished     bool      `json:"is_published"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Module struct {
	ID          int64     `json:"id"`
	CourseID    int64     `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OrderIndex  int       `json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
}

type Resource struct {
	ID            int64     `json:"id"`
	ModuleID      int64     `json:"module_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	ResourceType  string    `json:"resource_type"`
	Difficulty    int       `json:"difficulty"`
	EstimatedTime int       `json:"estimated_time"`
	FileURL       string    `json:"file_url"`
	ThumbnailURL  string    `json:"thumbnail_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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
