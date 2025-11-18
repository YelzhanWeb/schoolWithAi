package postgre

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"backend/internal/domain/models"
)

func (r *courseRepository) CreateModule(ctx context.Context, module *models.Module) error {
	query := `
		INSERT INTO modules (course_id, title, description, order_index, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx, query,
		module.CourseID,
		module.Title,
		module.Description,
		module.OrderIndex,
		time.Now(),
	).Scan(&module.ID)
	if err != nil {
		return fmt.Errorf("failed to create module: %w", err)
	}
	return nil
}

// UpdateModule обновляет модуль
func (r *courseRepository) UpdateModule(ctx context.Context, module *models.Module) error {
	query := `
		UPDATE modules
		SET title = $1, description = $2, order_index = $3
		WHERE id = $4 AND course_id = $5
	`
	result, err := r.db.ExecContext(
		ctx, query,
		module.Title,
		module.Description,
		module.OrderIndex,
		module.ID,
		module.CourseID,
	)
	if err != nil {
		return fmt.Errorf("failed to update module: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("module not found")
	}
	return nil
}

// DeleteModule удаляет модуль
func (r *courseRepository) DeleteModule(ctx context.Context, moduleID, courseID int64) error {
	query := `DELETE FROM modules WHERE id = $1 AND course_id = $2`

	result, err := r.db.ExecContext(ctx, query, moduleID, courseID)
	if err != nil {
		return fmt.Errorf("failed to delete module: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("module not found")
	}
	return nil
}

// CreateResource создает новый ресурс
func (r *courseRepository) CreateResource(ctx context.Context, resource *models.Resource) error {
	query := `
		INSERT INTO resources 
			(module_id, title, content, resource_type, difficulty, 
			 estimated_time, file_url, thumbnail_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx, query,
		resource.ModuleID,
		resource.Title,
		resource.Content,
		resource.ResourceType,
		resource.Difficulty,
		resource.EstimatedTime,
		resource.FileURL,
		resource.ThumbnailURL,
		time.Now(),
		time.Now(),
	).Scan(&resource.ID)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}
	return nil
}

// UpdateResource обновляет ресурс
func (r *courseRepository) UpdateResource(ctx context.Context, resource *models.Resource) error {
	query := `
		UPDATE resources
		SET title = $1, content = $2, resource_type = $3, difficulty = $4,
		    estimated_time = $5, file_url = $6, thumbnail_url = $7, updated_at = $8
		WHERE id = $9 AND module_id = $10
	`
	result, err := r.db.ExecContext(
		ctx, query,
		resource.Title,
		resource.Content,
		resource.ResourceType,
		resource.Difficulty,
		resource.EstimatedTime,
		resource.FileURL,
		resource.ThumbnailURL,
		time.Now(),
		resource.ID,
		resource.ModuleID,
	)
	if err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("resource not found")
	}
	return nil
}

// DeleteResource удаляет ресурс
func (r *courseRepository) DeleteResource(ctx context.Context, resourceID, moduleID int64) error {
	query := `DELETE FROM resources WHERE id = $1 AND module_id = $2`

	result, err := r.db.ExecContext(ctx, query, resourceID, moduleID)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("resource not found")
	}
	return nil
}

// PublishCourse публикует курс
func (r *courseRepository) PublishCourse(ctx context.Context, courseID, teacherID int64) error {
	query := `
		UPDATE courses
		SET is_published = true, updated_at = $1
		WHERE id = $2 AND created_by = $3
	`
	result, err := r.db.ExecContext(ctx, query, time.Now(), courseID, teacherID)
	if err != nil {
		return fmt.Errorf("failed to publish course: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("course not found or permission denied")
	}
	return nil
}

// UnpublishCourse снимает курс с публикации
func (r *courseRepository) UnpublishCourse(ctx context.Context, courseID, teacherID int64) error {
	query := `
		UPDATE courses
		SET is_published = false, updated_at = $1
		WHERE id = $2 AND created_by = $3
	`
	result, err := r.db.ExecContext(ctx, query, time.Now(), courseID, teacherID)
	if err != nil {
		return fmt.Errorf("failed to unpublish course: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("course not found or permission denied")
	}
	return nil
}

// DeleteCourse удаляет курс (каскадно удаляет модули и ресурсы)
func (r *courseRepository) Delete(ctx context.Context, courseID, teacherID int64) error {
	query := `DELETE FROM courses WHERE id = $1 AND created_by = $2`

	result, err := r.db.ExecContext(ctx, query, courseID, teacherID)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("course not found or permission denied")
	}
	return nil
}

// GetModuleByID получает модуль по ID с проверкой владельца курса
func (r *courseRepository) GetModuleByID(ctx context.Context, moduleID, teacherID int64) (*models.Module, error) {
	query := `
		SELECT m.id, m.course_id, m.title, m.description, m.order_index, m.created_at
		FROM modules m
		JOIN courses c ON m.course_id = c.id
		WHERE m.id = $1 AND c.created_by = $2
	`
	module := &models.Module{}
	err := r.db.QueryRowContext(ctx, query, moduleID, teacherID).Scan(
		&module.ID,
		&module.CourseID,
		&module.Title,
		&module.Description,
		&module.OrderIndex,
		&module.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("module not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get module: %w", err)
	}

	return module, nil
}
