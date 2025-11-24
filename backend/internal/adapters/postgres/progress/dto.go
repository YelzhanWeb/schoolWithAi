package progress

import (
	"time"

	"backend/internal/entities"
)

type lessonProgressDTO struct {
	UserID         string    `db:"user_id"`
	LessonID       string    `db:"lesson_id"`
	Status         string    `db:"status"`
	IsCompleted    bool      `db:"is_completed"`
	LastAccessedAt time.Time `db:"last_accessed_at"`
}

func (d *lessonProgressDTO) toEntity() (*entities.LessonProgress, error) {
	status, err := entities.ParseProgressStatus(d.Status)
	if err != nil {
		status = entities.StatusNotStarted
	}

	return &entities.LessonProgress{
		UserID:         d.UserID,
		LessonID:       d.LessonID,
		Status:         status,
		IsCompleted:    d.IsCompleted,
		LastAccessedAt: d.LastAccessedAt.UTC(),
	}, nil
}

type courseProgressDTO struct {
	UserID                string    `db:"user_id"`
	CourseID              string    `db:"course_id"`
	CompletedLessonsCount int       `db:"completed_lessons_count"`
	TotalLessonsCount     int       `db:"total_lessons_count"`
	ProgressPercentage    int       `db:"progress_percentage"`
	IsCompleted           bool      `db:"is_completed"`
	UpdatedAt             time.Time `db:"updated_at"`
}

func (d *courseProgressDTO) toEntity() *entities.CourseProgress {
	return &entities.CourseProgress{
		UserID:                d.UserID,
		CourseID:              d.CourseID,
		CompletedLessonsCount: d.CompletedLessonsCount,
		TotalLessonsCount:     d.TotalLessonsCount,
		ProgressPercentage:    d.ProgressPercentage,
		IsCompleted:           d.IsCompleted,
		UpdatedAt:             d.UpdatedAt.UTC(),
	}
}
