package repositories

import (
	"backend/internal/domain/models"
	"context"
)

type RecommendationRepository interface {
	// GetByStudent получить рекомендации студента
	GetByStudent(ctx context.Context, studentID int64, limit int) ([]*models.Recommendation, error)

	// Save сохранить рекомендации
	Save(ctx context.Context, recommendations []*models.Recommendation) error

	// MarkAsViewed отметить как просмотренное
	MarkAsViewed(ctx context.Context, recommendationID int64) error

	// DeleteOld удалить старые рекомендации
	DeleteOld(ctx context.Context, studentID int64, olderThanDays int) error
}
