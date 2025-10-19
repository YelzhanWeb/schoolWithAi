package repositories

import (
	"backend/internal/domain/models"
	"context"
)

// ProgressRepository - интерфейс для работы с прогрессом
type ProgressRepository interface {
	// Создать/обновить прогресс
	Upsert(ctx context.Context, progress *models.Progress) error

	// Получить прогресс студента по ресурсу
	GetByStudentAndResource(ctx context.Context, studentID, resourceID int64) (*models.Progress, error)

	// Получить весь прогресс студента
	GetByStudent(ctx context.Context, studentID int64) ([]*models.Progress, error)

	// Получить статистику студента
	GetStatistics(ctx context.Context, studentID int64) (*models.ProgressStats, error)
}
