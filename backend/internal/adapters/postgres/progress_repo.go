package postgre

import (
	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type progressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) repositories.ProgressRepository {
	return &progressRepository{db: db}
}

func (r *progressRepository) Upsert(ctx context.Context, progress *models.Progress) error {
	query := `
		INSERT INTO student_progress 
			(student_id, resource_id, status, score, time_spent, attempts, completed_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (student_id, resource_id) 
		DO UPDATE SET
			status = EXCLUDED.status,
			score = EXCLUDED.score,
			time_spent = student_progress.time_spent + EXCLUDED.time_spent,
			attempts = student_progress.attempts + 1,
			completed_at = EXCLUDED.completed_at,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		progress.StudentID,
		progress.ResourceID,
		progress.Status,
		progress.Score,
		progress.TimeSpent,
		progress.Attempts,
		progress.CompletedAt,
		time.Now(),
	).Scan(&progress.ID, &progress.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to upsert progress: %w", err)
	}

	return nil
}

func (r *progressRepository) GetByStudentAndResource(ctx context.Context, studentID, resourceID int64) (*models.Progress, error) {
	query := `
		SELECT id, student_id, resource_id, status, score, time_spent, attempts, 
		       completed_at, created_at, updated_at
		FROM student_progress
		WHERE student_id = $1 AND resource_id = $2
	`

	progress := &models.Progress{}
	err := r.db.QueryRowContext(ctx, query, studentID, resourceID).Scan(
		&progress.ID,
		&progress.StudentID,
		&progress.ResourceID,
		&progress.Status,
		&progress.Score,
		&progress.TimeSpent,
		&progress.Attempts,
		&progress.CompletedAt,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // не найдено - нормально
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	return progress, nil
}

func (r *progressRepository) GetByStudent(ctx context.Context, studentID int64) ([]*models.Progress, error) {
	query := `
		SELECT p.id, p.student_id, p.resource_id, p.status, p.score, p.time_spent, 
		       p.attempts, p.completed_at, p.created_at, p.updated_at
		FROM student_progress p
		WHERE p.student_id = $1
		ORDER BY p.updated_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query progress: %w", err)
	}
	defer rows.Close()

	var progressList []*models.Progress
	for rows.Next() {
		progress := &models.Progress{}
		err := rows.Scan(
			&progress.ID,
			&progress.StudentID,
			&progress.ResourceID,
			&progress.Status,
			&progress.Score,
			&progress.TimeSpent,
			&progress.Attempts,
			&progress.CompletedAt,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan progress: %w", err)
		}
		progressList = append(progressList, progress)
	}

	return progressList, nil
}

func (r *progressRepository) GetStatistics(ctx context.Context, studentID int64) (*models.ProgressStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_resources,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_count,
			SUM(CASE WHEN status = 'in_progress' THEN 1 ELSE 0 END) as in_progress_count,
			COALESCE(AVG(CASE WHEN score > 0 THEN score END), 0) as avg_score,
			COALESCE(SUM(time_spent), 0) as total_time_spent
		FROM student_progress
		WHERE student_id = $1
	`

	stats := &models.ProgressStats{}
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(
		&stats.TotalResources,
		&stats.CompletedCount,
		&stats.InProgressCount,
		&stats.AverageScore,
		&stats.TotalTimeSpent,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	// Рассчитать процент завершения
	if stats.TotalResources > 0 {
		stats.CompletionRate = float64(stats.CompletedCount) / float64(stats.TotalResources) * 100
	}

	// Получить очки и уровень из student_points
	pointsQuery := `
		SELECT COALESCE(total_points, 0), COALESCE(level, 1)
		FROM student_points
		WHERE student_id = $1
	`
	_ = r.db.QueryRowContext(ctx, pointsQuery, studentID).Scan(&stats.TotalPoints, &stats.Level)

	return stats, nil
}
