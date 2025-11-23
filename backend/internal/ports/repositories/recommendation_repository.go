package repositories

import (
	"context"

	models "backend/internal/entities"
)

// CourseRecommendationRepository определяет интерфейс
type CourseRecommendationRepository interface {
	// GetByStudent получить рекомендации курсов студента
	GetByStudent(ctx context.Context, studentID int64, limit int) ([]*models.CourseRecommendation, error)

	// Save сохранить рекомендации курсов
	Save(ctx context.Context, recommendations []*models.CourseRecommendation) error

	// MarkAsViewed отметить как просмотренное
	MarkAsViewed(ctx context.Context, recommendationID int64) error

	// DeleteOld удалить старые рекомендации курсов
	DeleteOld(ctx context.Context, studentID int64, olderThanDays int) error
}
