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

func NewRecommendationRepository(db *sql.DB) repositories.RecommendationRepository {
	return &recommendationRepository{db: db}
}

func (r *recommendationRepository) GetByStudent(ctx context.Context, studentID int64, limit int) ([]*models.Recommendation, error) {
	query := `
		SELECT 
			r.id, r.student_id, r.resource_id, r.score, 
			r.reason, r.algorithm_type, r.is_viewed, r.created_at,
			res.title
		FROM recommendations r
		JOIN resources res ON r.resource_id = res.id
		WHERE r.student_id = $1
		ORDER BY r.score DESC, r.created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, studentID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recommendations: %w", err)
	}
	defer rows.Close()

	var recommendations []*models.Recommendation
	for rows.Next() {
		rec := &models.Recommendation{}
		var title string

		err := rows.Scan(
			&rec.ID,
			&rec.StudentID,
			&rec.ResourceID,
			&rec.Score,
			&rec.Reason,
			&rec.AlgorithmType,
			&rec.IsViewed,
			&rec.CreatedAt,
			&title,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recommendation: %w", err)
		}

		recommendations = append(recommendations, rec)
	}

	return recommendations, nil
}

func (r *recommendationRepository) Save(ctx context.Context, recommendations []*models.Recommendation) error {
	query := `
		INSERT INTO recommendations 
			(student_id, resource_id, score, reason, algorithm_type, is_viewed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (student_id, resource_id) DO UPDATE
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
			rec.ResourceID,
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
	query := `UPDATE recommendations SET is_viewed = true WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, recommendationID)
	if err != nil {
		return fmt.Errorf("failed to mark as viewed: %w", err)
	}

	return nil
}

func (r *recommendationRepository) DeleteOld(ctx context.Context, studentID int64, olderThanDays int) error {
	query := `
		DELETE FROM recommendations
		WHERE student_id = $1 
		  AND created_at < NOW() - INTERVAL '%d days'
	`

	_, err := r.db.ExecContext(ctx, fmt.Sprintf(query, olderThanDays), studentID)
	if err != nil {
		return fmt.Errorf("failed to delete old recommendations: %w", err)
	}

	return nil
}
