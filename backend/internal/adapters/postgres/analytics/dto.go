package analytics

import (
	"time"

	"backend/internal/entities"
)

type dto struct {
	ID         string         `db:"id"`
	UserID     string         `db:"user_id"`
	CourseID   *string        `db:"course_id"` // Указатель для NULL
	ActionType string         `db:"action_type"`
	MetaData   map[string]any `db:"meta_data"`
	CreatedAt  time.Time      `db:"created_at"`
}

func newDTO(e *entities.UserActivityLog) dto {
	return dto{
		ID:         e.ID,
		UserID:     e.UserID,
		CourseID:   e.CourseID,
		ActionType: e.ActionType,
		MetaData:   e.MetaData,
		CreatedAt:  e.CreatedAt,
	}
}

func (d *dto) toEntity() *entities.UserActivityLog {
	return &entities.UserActivityLog{
		ID:         d.ID,
		UserID:     d.UserID,
		CourseID:   d.CourseID,
		ActionType: d.ActionType,
		MetaData:   d.MetaData,
		CreatedAt:  d.CreatedAt.UTC(),
	}
}
