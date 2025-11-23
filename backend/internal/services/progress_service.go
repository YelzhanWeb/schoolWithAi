package services

import (
	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
	"context"
	"errors"
	"time"
)

type ProgressService struct {
	progressRepo repositories.ProgressRepository
}

func NewProgressService(progressRepo repositories.ProgressRepository) *ProgressService {
	return &ProgressService{
		progressRepo: progressRepo,
	}
}

// UpdateProgress - обновить прогресс урока
func (s *ProgressService) UpdateProgress(ctx context.Context, update *models.ProgressUpdate) error {
	// Валидация
	if update.Status != "not_started" && update.Status != "in_progress" && update.Status != "completed" {
		return errors.New("invalid status")
	}

	if update.Score < 0 || update.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}

	// Получить текущий прогресс
	existing, err := s.progressRepo.GetByStudentAndResource(ctx, update.StudentID, update.ResourceID)
	if err != nil {
		return err
	}

	progress := &models.Progress{
		StudentID:  update.StudentID,
		ResourceID: update.ResourceID,
		Status:     update.Status,
		Score:      update.Score,
		TimeSpent:  update.TimeSpent,
		Attempts:   1,
	}

	// Если завершено - установить время
	if update.Status == "completed" {
		now := time.Now()
		progress.CompletedAt = &now
	}

	// Если уже есть прогресс - обновляем
	if existing != nil {
		progress.ID = existing.ID
		progress.Attempts = existing.Attempts + 1
	}

	return s.progressRepo.Upsert(ctx, progress)
}

// GetStudentProgress - получить прогресс студента
func (s *ProgressService) GetStudentProgress(ctx context.Context, studentID int64) ([]*models.Progress, error) {
	return s.progressRepo.GetByStudent(ctx, studentID)
}

// GetStudentStatistics - получить статистику студента
func (s *ProgressService) GetStudentStatistics(ctx context.Context, studentID int64) (*models.ProgressStats, error) {
	return s.progressRepo.GetStatistics(ctx, studentID)
}
