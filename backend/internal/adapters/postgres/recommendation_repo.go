package postgre

import (
	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type recommendationRepository struct {
	db *sql.DB
}

func NewRecommendationRepository(db *sql.DB) repositories.CourseRecommendationRepository {
	return &recommendationRepository{db: db}
}

func (r *recommendationRepository) GetByStudent(ctx context.Context, studentID int64, limit int) ([]*models.CourseRecommendation, error) {
	query := `
		SELECT
    cr.id, cr.student_id, cr.course_id, cr.score,
    cr.reason, cr.algorithm_type, cr.is_viewed, cr.created_at,
    c.title -- Добавлено поле title из таблицы courses
FROM course_recommendations cr -- Изменено
JOIN courses c ON cr.course_id = c.id -- Добавлен JOIN
WHERE cr.student_id = $1
ORDER BY cr.score DESC, cr.created_at DESC
LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, studentID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recommendations: %w", err)
	}
	defer rows.Close()

	var recommendations []*models.CourseRecommendation
	for rows.Next() {
		rec := &models.CourseRecommendation{}

		err := rows.Scan(
			&rec.ID,
			&rec.StudentID,
			&rec.CourseID,
			&rec.Score,
			&rec.Reason,
			&rec.AlgorithmType,
			&rec.IsViewed,
			&rec.CreatedAt,
			&rec.CourseTitle,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recommendation: %w", err)
		}

		recommendations = append(recommendations, rec)
	}

	return recommendations, nil
}

func (r *recommendationRepository) Save(ctx context.Context, recommendations []*models.CourseRecommendation) error {
	query := `
		INSERT INTO course_recommendations -- Изменено
    (student_id, course_id, score, reason, algorithm_type, is_viewed, created_at) -- Изменено
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (student_id, course_id) DO UPDATE -- Изменено
SET score = EXCLUDED.score,
    reason = EXCLUDED.reason,
    algorithm_type = EXCLUDED.algorithm_type,
    created_at = EXCLUDED.created_at
	`

	for _, rec := range recommendations {
		_, err := r.db.ExecContext(
			ctx,
			query,
			rec.StudentID,
			rec.CourseID,
			rec.Score,
			rec.Reason,
			rec.AlgorithmType,
			false,
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to save recommendation: %w", err)
		}
	}

	return nil
}

func (r *recommendationRepository) MarkAsViewed(ctx context.Context, recommendationID int64) error {
	query := `UPDATE course_recommendations SET is_viewed = true WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, recommendationID)
	if err != nil {
		return fmt.Errorf("failed to mark as viewed: %w", err)
	}

	return nil
}

func (r *recommendationRepository) DeleteOld(ctx context.Context, studentID int64, olderThanDays int) error {
	query := `
		DELETE FROM course_recommendations -- Изменено
WHERE student_id = $1
  AND created_at < NOW() - INTERVAL '%d days'
	`

	_, err := r.db.ExecContext(ctx, fmt.Sprintf(query, olderThanDays), studentID)
	if err != nil {
		return fmt.Errorf("failed to delete old recommendations: %w", err)
	}

	return nil
}
