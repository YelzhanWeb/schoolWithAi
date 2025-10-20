package postgre

import (
	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
	"context"
	"database/sql"
	"fmt"
)

type courseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) repositories.CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) GetAll(ctx context.Context) ([]*models.Course, error) {
	query := `
		SELECT 
			c.id, 
			c.title, 
			c.description, 
			c.created_by,
			c.difficulty_level, 
			c.age_group,
			c.subject,
			c.is_published,
			c.thumbnail_url,
			c.created_at,
			c.updated_at
		FROM courses c
		WHERE c.is_published = true
		ORDER BY c.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query courses: %w", err)
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		course := &models.Course{}
		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.CreatedBy,
			&course.DifficultyLevel,
			&course.AgeGroup,
			&course.Subject,
			&course.IsPublished,
			&course.ThumbnailURL,
			&course.CreatedAt,
			&course.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (r *courseRepository) GetByID(ctx context.Context, id int64) (*models.Course, error) {
	query := `
		SELECT 
			id, title, description, created_by, difficulty_level, 
			age_group, subject, is_published, thumbnail_url, 
			created_at, updated_at
		FROM courses
		WHERE id = $1 AND is_published = true
	`

	course := &models.Course{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.CreatedBy,
		&course.DifficultyLevel,
		&course.AgeGroup,
		&course.Subject,
		&course.IsPublished,
		&course.ThumbnailURL,
		&course.CreatedAt,
		&course.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("course not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	return course, nil
}

func (r *courseRepository) GetModules(ctx context.Context, courseID int64) ([]*models.Module, error) {
	query := `
		SELECT id, course_id, title, description, order_index, created_at
		FROM modules
		WHERE course_id = $1
		ORDER BY order_index
	`

	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query modules: %w", err)
	}
	defer rows.Close()

	var modules []*models.Module
	for rows.Next() {
		module := &models.Module{}
		err := rows.Scan(
			&module.ID,
			&module.CourseID,
			&module.Title,
			&module.Description,
			&module.OrderIndex,
			&module.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		modules = append(modules, module)
	}

	return modules, nil
}

func (r *courseRepository) GetResources(ctx context.Context, moduleID int64) ([]*models.Resource, error) {
	query := `
		SELECT 
			id, module_id, title, content, resource_type, 
			difficulty, estimated_time, file_url, thumbnail_url,
			created_at, updated_at
		FROM resources
		WHERE module_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to query resources: %w", err)
	}
	defer rows.Close()

	var resources []*models.Resource
	for rows.Next() {
		resource := &models.Resource{}
		err := rows.Scan(
			&resource.ID,
			&resource.ModuleID,
			&resource.Title,
			&resource.Content,
			&resource.ResourceType,
			&resource.Difficulty,
			&resource.EstimatedTime,
			&resource.FileURL,
			&resource.ThumbnailURL,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan resource: %w", err)
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

func (r *courseRepository) Create(ctx context.Context, course *models.Course) error {
	// TODO: Implement
	return nil
}

func (r *courseRepository) Update(ctx context.Context, course *models.Course) error {
	// TODO: Implement
	return nil
}

func (r *courseRepository) GetByTeacher(ctx context.Context, teacherID int64) ([]*models.Course, error) {
	// TODO: Implement
	return nil, nil
}

// backend/internal/adapters/postgres/course_repo.go
func (r *courseRepository) GetResourceByID(ctx context.Context, id int64) (*models.Resource, error) {
	query := `
		SELECT 
			id, module_id, title, content, resource_type, 
			difficulty, estimated_time, file_url, thumbnail_url,
			created_at, updated_at
		FROM resources
		WHERE id = $1
	`

	resource := &models.Resource{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&resource.ID,
		&resource.ModuleID,
		&resource.Title,
		&resource.Content,
		&resource.ResourceType,
		&resource.Difficulty,
		&resource.EstimatedTime,
		&resource.FileURL,
		&resource.ThumbnailURL,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("resource not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return resource, nil
}
