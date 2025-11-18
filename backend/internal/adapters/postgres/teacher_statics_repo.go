package postgre

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
)

type teacherStatisticsRepository struct {
	db *sql.DB
}

func NewTeacherStatisticsRepository(db *sql.DB) repositories.TeacherStatisticsRepository {
	return &teacherStatisticsRepository{db: db}
}

// GetTeacherStatistics получает общую статистику учителя
func (r *teacherStatisticsRepository) GetTeacherStatistics(ctx context.Context, teacherID int64) (*models.TeacherStatistics, error) {
	stats := &models.TeacherStatistics{}

	// Общая статистика по курсам
	query := `
		SELECT 
			COUNT(*) as total_courses,
			SUM(CASE WHEN is_published = true THEN 1 ELSE 0 END) as published_courses,
			SUM(CASE WHEN is_published = false THEN 1 ELSE 0 END) as draft_courses
		FROM courses
		WHERE created_by = $1
	`
	err := r.db.QueryRowContext(ctx, query, teacherID).Scan(
		&stats.TotalCourses,
		&stats.PublishedCourses,
		&stats.DraftCourses,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get course stats: %w", err)
	}

	// Общее количество студентов (уникальных)
	studentQuery := `
		SELECT COUNT(DISTINCT sp.student_id)
		FROM student_progress sp
		JOIN resources r ON sp.resource_id = r.id
		JOIN modules m ON r.module_id = m.id
		JOIN courses c ON m.course_id = c.id
		WHERE c.created_by = $1
	`
	err = r.db.QueryRowContext(ctx, studentQuery, teacherID).Scan(&stats.TotalStudents)
	if err != nil {
		stats.TotalStudents = 0
	}

	// Количество модулей
	moduleQuery := `
		SELECT COUNT(m.id)
		FROM modules m
		JOIN courses c ON m.course_id = c.id
		WHERE c.created_by = $1
	`
	err = r.db.QueryRowContext(ctx, moduleQuery, teacherID).Scan(&stats.TotalModules)
	if err != nil {
		stats.TotalModules = 0
	}

	// Количество ресурсов
	resourceQuery := `
		SELECT COUNT(r.id)
		FROM resources r
		JOIN modules m ON r.module_id = m.id
		JOIN courses c ON m.course_id = c.id
		WHERE c.created_by = $1
	`
	err = r.db.QueryRowContext(ctx, resourceQuery, teacherID).Scan(&stats.TotalResources)
	if err != nil {
		stats.TotalResources = 0
	}

	return stats, nil
}

// GetCourseStatistics получает статистику по конкретному курсу
func (r *teacherStatisticsRepository) GetCourseStatistics(ctx context.Context, courseID, teacherID int64) (*models.CourseStatistics, error) {
	stats := &models.CourseStatistics{
		CourseID: courseID,
	}

	// Название курса
	query := `
		SELECT title FROM courses 
		WHERE id = $1 AND created_by = $2
	`
	err := r.db.QueryRowContext(ctx, query, courseID, teacherID).Scan(&stats.CourseTitle)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	// Статистика по студентам
	statsQuery := `
		SELECT 
			COUNT(DISTINCT sp.student_id) as enrolled_students,
			SUM(CASE WHEN sp.status = 'completed' THEN 1 ELSE 0 END) as completed_count,
			SUM(CASE WHEN sp.status = 'in_progress' THEN 1 ELSE 0 END) as in_progress_count,
			COALESCE(AVG(CASE WHEN sp.score > 0 THEN sp.score END), 0) as average_score
		FROM student_progress sp
		JOIN resources r ON sp.resource_id = r.id
		JOIN modules m ON r.module_id = m.id
		WHERE m.course_id = $1
	`
	err = r.db.QueryRowContext(ctx, statsQuery, courseID).Scan(
		&stats.EnrolledStudents,
		&stats.CompletedCount,
		&stats.InProgressCount,
		&stats.AverageScore,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	// Рассчитываем процент завершения
	totalInteractions := stats.CompletedCount + stats.InProgressCount
	if totalInteractions > 0 {
		stats.CompletionRate = float64(stats.CompletedCount) / float64(totalInteractions) * 100
	}

	// Средний рейтинг курса
	ratingQuery := `
		SELECT 
			COALESCE(AVG(rr.rating), 0) as average_rating,
			COUNT(rr.id) as total_ratings
		FROM resource_ratings rr
		JOIN resources r ON rr.resource_id = r.id
		JOIN modules m ON r.module_id = m.id
		WHERE m.course_id = $1
	`
	err = r.db.QueryRowContext(ctx, ratingQuery, courseID).Scan(
		&stats.AverageRating,
		&stats.TotalRatings,
	)
	if err != nil {
		stats.AverageRating = 0
		stats.TotalRatings = 0
	}

	return stats, nil
}

// GetAllCourseStatistics получает статистику по всем курсам учителя
func (r *teacherStatisticsRepository) GetAllCourseStatistics(ctx context.Context, teacherID int64) ([]*models.CourseStatistics, error) {
	query := `
		SELECT c.id
		FROM courses c
		WHERE c.created_by = $1
		ORDER BY c.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to query courses: %w", err)
	}
	defer rows.Close()

	var courseStats []*models.CourseStatistics
	for rows.Next() {
		var courseID int64
		if err := rows.Scan(&courseID); err != nil {
			continue
		}

		stats, err := r.GetCourseStatistics(ctx, courseID, teacherID)
		if err != nil {
			continue
		}
		courseStats = append(courseStats, stats)
	}

	return courseStats, nil
}

// GetStudentProgressByCourse получает детальный прогресс студентов по курсу
func (r *teacherStatisticsRepository) GetStudentProgressByCourse(ctx context.Context, courseID, teacherID int64) ([]*models.StudentCourseProgress, error) {
	// Проверка прав
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1 AND created_by = $2)`
	err := r.db.QueryRowContext(ctx, checkQuery, courseID, teacherID).Scan(&exists)
	if err != nil || !exists {
		return nil, fmt.Errorf("course not found or permission denied")
	}

	query := `
		SELECT 
			u.id as student_id,
			u.full_name as student_name,
			COUNT(DISTINCT r.id) as total_resources,
			SUM(CASE WHEN sp.status = 'completed' THEN 1 ELSE 0 END) as completed_count,
			SUM(CASE WHEN sp.status = 'in_progress' THEN 1 ELSE 0 END) as in_progress_count,
			COALESCE(AVG(CASE WHEN sp.score > 0 THEN sp.score END), 0) as average_score,
			COALESCE(SUM(sp.time_spent), 0) as total_time_spent,
			MAX(sp.updated_at) as last_activity_date
		FROM users u
		JOIN student_progress sp ON u.id = sp.student_id
		JOIN resources r ON sp.resource_id = r.id
		JOIN modules m ON r.module_id = m.id
		WHERE m.course_id = $1 AND u.role = 'student'
		GROUP BY u.id, u.full_name
		ORDER BY last_activity_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query student progress: %w", err)
	}
	defer rows.Close()

	var progressList []*models.StudentCourseProgress
	for rows.Next() {
		progress := &models.StudentCourseProgress{
			CourseID: courseID,
		}
		err := rows.Scan(
			&progress.StudentID,
			&progress.StudentName,
			&progress.TotalResources,
			&progress.CompletedCount,
			&progress.InProgressCount,
			&progress.AverageScore,
			&progress.TotalTimeSpent,
			&progress.LastActivityDate,
		)
		if err != nil {
			continue
		}

		// Рассчитываем процент завершения
		if progress.TotalResources > 0 {
			progress.CompletionRate = float64(progress.CompletedCount) / float64(progress.TotalResources) * 100
		}

		progressList = append(progressList, progress)
	}

	return progressList, nil
}
