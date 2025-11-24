package entities

import (
	"errors"
	"time"
)

type ProgressStatus string

const (
	StatusNotStarted ProgressStatus = "not_started"
	StatusInProgress ProgressStatus = "in_progress"
	StatusCompleted  ProgressStatus = "completed"
)

type LessonProgress struct {
	UserID         string
	LessonID       string
	Status         ProgressStatus
	IsCompleted    bool
	LastAccessedAt time.Time
}

type CourseProgress struct {
	UserID                string
	CourseID              string
	CompletedLessonsCount int
	TotalLessonsCount     int
	ProgressPercentage    int // 0-100
	IsCompleted           bool
	UpdatedAt             time.Time
}

func NewLessonProgress(userID, lessonID string) *LessonProgress {
	return &LessonProgress{
		UserID:         userID,
		LessonID:       lessonID,
		Status:         StatusNotStarted,
		IsCompleted:    false,
		LastAccessedAt: time.Now().UTC(),
	}
}

func (s ProgressStatus) IsValid() bool {
	switch s {
	case StatusNotStarted, StatusInProgress, StatusCompleted:
		return true
	}
	return false
}

func ParseProgressStatus(s string) (ProgressStatus, error) {
	status := ProgressStatus(s)
	if !status.IsValid() {
		return "", errors.New("invalid progress status")
	}
	return status, nil
}
