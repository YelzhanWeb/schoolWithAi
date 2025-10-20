package repositories

import (
	"backend/internal/domain/models"
	"context"
)

// CourseRepository определяет контракт для работы с курсами
type CourseRepository interface {
	// GetAll возвращает все опубликованные курсы
	GetAll(ctx context.Context) ([]*models.Course, error)

	// GetByID возвращает курс по ID
	GetByID(ctx context.Context, id int64) (*models.Course, error)

	// GetModules возвращает модули курса
	GetModules(ctx context.Context, courseID int64) ([]*models.Module, error)

	// GetResources возвращает ресурсы модуля
	GetResources(ctx context.Context, moduleID int64) ([]*models.Resource, error)

	// Create создает новый курс (для учителя)
	Create(ctx context.Context, course *models.Course) error

	// Update обновляет курс
	Update(ctx context.Context, course *models.Course) error

	// GetByTeacher возвращает курсы учителя
	GetByTeacher(ctx context.Context, teacherID int64) ([]*models.Course, error)

	// backend/internal/ports/repositories/course_repository.go
	GetResourceByID(ctx context.Context, id int64) (*models.Resource, error)
}
