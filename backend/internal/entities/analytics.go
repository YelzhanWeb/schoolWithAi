package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserActivityLog struct {
	ID         string
	UserID     string
	CourseID   *string        // Может быть nil, если действие не привязано к курсу
	ActionType string         // 'view', 'complete', 'search', 'like'
	MetaData   map[string]any // Гибкие данные: {"duration": 120, "query": "python"}
	CreatedAt  time.Time
}

// NewActivityLog создает лог. courseID передаем как *string (или nil).
func NewActivityLog(
	userID string,
	courseID *string,
	actionType string,
	metaData map[string]any,
) (*UserActivityLog, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	if actionType == "" {
		return nil, errors.New("action_type is required")
	}

	// Если мета-данных нет, создаем пустую мапу, чтобы в базу записалось "{}"
	if metaData == nil {
		metaData = make(map[string]any)
	}

	return &UserActivityLog{
		ID:         uuid.NewString(),
		UserID:     userID,
		CourseID:   courseID,
		ActionType: actionType,
		MetaData:   metaData,
		CreatedAt:  time.Now().UTC(),
	}, nil
}
