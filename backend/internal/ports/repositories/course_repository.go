package repositories

import (
	"context"

	"backend/internal/domain/models"
)

// CourseRepository - расширенный интерфейс с методами для учителя
type CourseRepository interface {
	// Общие методы (уже есть)
	GetAll(ctx context.Context) ([]*models.Course, error)
	GetByID(ctx context.Context, id int64) (*models.Course, error)
	GetModules(ctx context.Context, courseID int64) ([]*models.Module, error)
	GetResources(ctx context.Context, moduleID int64) ([]*models.Resource, error)
	GetResourceByID(ctx context.Context, id int64) (*models.Resource, error)
	GetByTeacher(ctx context.Context, teacherID int64) ([]*models.Course, error)

	// Управление курсами
	Create(ctx context.Context, course *models.Course) error
	Update(ctx context.Context, course *models.Course) error
	Delete(ctx context.Context, courseID, teacherID int64) error
	PublishCourse(ctx context.Context, courseID, teacherID int64) error
	UnpublishCourse(ctx context.Context, courseID, teacherID int64) error

	// Управление модулями
	CreateModule(ctx context.Context, module *models.Module) error
	UpdateModule(ctx context.Context, module *models.Module) error
	DeleteModule(ctx context.Context, moduleID, courseID int64) error
	GetModuleByID(ctx context.Context, moduleID, teacherID int64) (*models.Module, error)

	// Управление ресурсами
	CreateResource(ctx context.Context, resource *models.Resource) error
	UpdateResource(ctx context.Context, resource *models.Resource) error
	DeleteResource(ctx context.Context, resourceID, moduleID int64) error
}

// TeacherStatisticsRepository - репозиторий статистики учителя
type TeacherStatisticsRepository interface {
	// Общая статистика учителя
	GetTeacherStatistics(ctx context.Context, teacherID int64) (*models.TeacherStatistics, error)

	// Статистика по курсам
	GetCourseStatistics(ctx context.Context, courseID, teacherID int64) (*models.CourseStatistics, error)
	GetAllCourseStatistics(ctx context.Context, teacherID int64) ([]*models.CourseStatistics, error)

	// Прогресс студентов
	GetStudentProgressByCourse(ctx context.Context, courseID, teacherID int64) ([]*models.StudentCourseProgress, error)
}
