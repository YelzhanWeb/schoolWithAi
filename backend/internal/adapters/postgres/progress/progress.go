package progress

import (
	"context"
	"fmt"

	"backend/internal/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProgressRepository struct {
	pool *pgxpool.Pool
}

func NewProgressRepository(pool *pgxpool.Pool) *ProgressRepository {
	return &ProgressRepository{pool: pool}
}

func (r *ProgressRepository) UpsertLessonProgress(ctx context.Context, lp *entities.LessonProgress) error {
	query := `
		INSERT INTO lesson_progress (user_id, lesson_id, status, is_completed, last_accessed_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, lesson_id) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			is_completed = EXCLUDED.is_completed,
			last_accessed_at = EXCLUDED.last_accessed_at
	`

	_, err := r.pool.Exec(ctx, query, lp.UserID, lp.LessonID, string(lp.Status), lp.IsCompleted, lp.LastAccessedAt)
	if err != nil {
		return fmt.Errorf("upsert lesson progress: %w", err)
	}
	return nil
}

func (r *ProgressRepository) GetLessonProgress(ctx context.Context, userID, lessonID string) (*entities.LessonProgress, error) {
	query := `
		SELECT user_id, lesson_id, status, is_completed, last_accessed_at
		FROM lesson_progress
		WHERE user_id = $1 AND lesson_id = $2
	`

	var d lessonProgressDTO
	err := r.pool.QueryRow(ctx, query, userID, lessonID).Scan(
		&d.UserID, &d.LessonID, &d.Status, &d.IsCompleted, &d.LastAccessedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.NewLessonProgress(userID, lessonID), nil
		}
		return nil, fmt.Errorf("get lesson progress: %w", err)
	}

	return d.toEntity()
}

func (r *ProgressRepository) CountCompletedLessonsInCourse(ctx context.Context, userID, courseID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM lesson_progress lp
		JOIN lessons l ON lp.lesson_id = l.id
		JOIN modules m ON l.module_id = m.id
		WHERE lp.user_id = $1 
		  AND m.course_id = $2 
		  AND lp.is_completed = true
	`

	var count int
	err := r.pool.QueryRow(ctx, query, userID, courseID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count completed lessons: %w", err)
	}
	return count, nil
}

func (r *ProgressRepository) UpsertCourseProgress(ctx context.Context, cp *entities.CourseProgress) error {
	query := `
		INSERT INTO course_progress (
			user_id, course_id, 
			completed_lessons_count, total_lessons_count, progress_percentage, 
			is_completed, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, course_id) 
		DO UPDATE SET 
			completed_lessons_count = EXCLUDED.completed_lessons_count,
			total_lessons_count = EXCLUDED.total_lessons_count,
			progress_percentage = EXCLUDED.progress_percentage,
			is_completed = EXCLUDED.is_completed,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.pool.Exec(
		ctx, query,
		cp.UserID, cp.CourseID,
		cp.CompletedLessonsCount, cp.TotalLessonsCount, cp.ProgressPercentage,
		cp.IsCompleted, cp.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert course progress: %w", err)
	}
	return nil
}

func (r *ProgressRepository) GetCourseProgress(ctx context.Context, userID, courseID string) (*entities.CourseProgress, error) {
	query := `
		SELECT user_id, course_id, completed_lessons_count, total_lessons_count, 
		       progress_percentage, is_completed, updated_at
		FROM course_progress
		WHERE user_id = $1 AND course_id = $2
	`

	var d courseProgressDTO
	err := r.pool.QueryRow(ctx, query, userID, courseID).Scan(
		&d.UserID, &d.CourseID,
		&d.CompletedLessonsCount, &d.TotalLessonsCount,
		&d.ProgressPercentage, &d.IsCompleted, &d.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get course progress: %w", err)
	}

	return d.toEntity(), nil
}

func (r *ProgressRepository) GetUserActiveCourses(ctx context.Context, userID string, limit int) ([]entities.CourseProgress, error) {
	query := `
		SELECT user_id, course_id, completed_lessons_count, total_lessons_count, 
		       progress_percentage, is_completed, updated_at
		FROM course_progress
		WHERE user_id = $1 AND is_completed = false
		ORDER BY updated_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get active courses: %w", err)
	}
	defer rows.Close()

	var list []entities.CourseProgress
	for rows.Next() {
		var d courseProgressDTO
		err := rows.Scan(
			&d.UserID, &d.CourseID,
			&d.CompletedLessonsCount, &d.TotalLessonsCount,
			&d.ProgressPercentage, &d.IsCompleted, &d.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, *d.toEntity())
	}

	return list, nil
}
