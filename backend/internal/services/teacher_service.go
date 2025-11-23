package services

import (
	"context"
	"errors"
	"fmt"

	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
)

type TeacherService struct {
	courseRepo repositories.CourseRepository
	statsRepo  repositories.TeacherStatisticsRepository
}

func NewTeacherService(
	courseRepo repositories.CourseRepository,
	statsRepo repositories.TeacherStatisticsRepository,
) *TeacherService {
	return &TeacherService{
		courseRepo: courseRepo,
		statsRepo:  statsRepo,
	}
}

// ========================================
// УПРАВЛЕНИЕ МОДУЛЯМИ
// ========================================

func (s *TeacherService) CreateModule(ctx context.Context, module *models.Module, teacherID int64) error {
	// Проверяем, что курс принадлежит учителю
	course, err := s.courseRepo.GetByID(ctx, module.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	if course.CreatedBy != teacherID {
		return errors.New("permission denied: you don't own this course")
	}

	// Валидация
	if module.Title == "" {
		return errors.New("module title is required")
	}

	if module.OrderIndex < 0 {
		return errors.New("order index must be non-negative")
	}

	return s.courseRepo.CreateModule(ctx, module)
}

func (s *TeacherService) UpdateModule(ctx context.Context, module *models.Module, teacherID int64) error {
	// Проверяем права через GetModuleByID
	_, err := s.courseRepo.GetModuleByID(ctx, module.ID, teacherID)
	if err != nil {
		return errors.New("module not found or permission denied")
	}

	if module.Title == "" {
		return errors.New("module title is required")
	}

	return s.courseRepo.UpdateModule(ctx, module)
}

func (s *TeacherService) DeleteModule(ctx context.Context, moduleID, courseID, teacherID int64) error {
	// Проверяем права
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return errors.New("course not found")
	}

	if course.CreatedBy != teacherID {
		return errors.New("permission denied")
	}

	return s.courseRepo.DeleteModule(ctx, moduleID, courseID)
}

func (s *TeacherService) GetModulesByTeacher(ctx context.Context, moduleID, teacherID int64) (*models.Module, error) {
	return s.courseRepo.GetModuleByID(ctx, moduleID, teacherID)
}

// ========================================
// УПРАВЛЕНИЕ РЕСУРСАМИ
// ========================================

func (s *TeacherService) CreateResource(ctx context.Context, resource *models.Resource, teacherID int64) error {
	// Проверяем права через модуль
	module, err := s.courseRepo.GetModuleByID(ctx, resource.ModuleID, teacherID)
	if err != nil {
		return errors.New("module not found or permission denied")
	}

	// Валидация
	if resource.Title == "" {
		return errors.New("resource title is required")
	}

	if resource.Difficulty < 1 || resource.Difficulty > 5 {
		return errors.New("difficulty must be between 1 and 5")
	}

	validTypes := map[string]bool{
		"exercise":    true,
		"quiz":        true,
		"reading":     true,
		"video":       true,
		"interactive": true,
	}
	if !validTypes[resource.ResourceType] {
		return errors.New("invalid resource type")
	}

	// Устанавливаем ModuleID из проверенного модуля
	resource.ModuleID = module.ID

	return s.courseRepo.CreateResource(ctx, resource)
}

func (s *TeacherService) UpdateResource(ctx context.Context, resource *models.Resource, teacherID int64) error {
	// Проверяем права через модуль
	_, err := s.courseRepo.GetModuleByID(ctx, resource.ModuleID, teacherID)
	if err != nil {
		return errors.New("permission denied")
	}

	if resource.Title == "" {
		return errors.New("resource title is required")
	}

	if resource.Difficulty < 1 || resource.Difficulty > 5 {
		return errors.New("difficulty must be between 1 and 5")
	}

	return s.courseRepo.UpdateResource(ctx, resource)
}

func (s *TeacherService) DeleteResource(ctx context.Context, resourceID, moduleID, teacherID int64) error {
	// Проверяем права
	_, err := s.courseRepo.GetModuleByID(ctx, moduleID, teacherID)
	if err != nil {
		return errors.New("permission denied")
	}

	return s.courseRepo.DeleteResource(ctx, resourceID, moduleID)
}

// ========================================
// УПРАВЛЕНИЕ КУРСАМИ
// ========================================

func (s *TeacherService) PublishCourse(ctx context.Context, courseID, teacherID int64) error {
	// Проверяем, что курс готов к публикации
	modules, err := s.courseRepo.GetModules(ctx, courseID)
	if err != nil {
		return errors.New("failed to check course modules")
	}

	if len(modules) == 0 {
		return errors.New("cannot publish course without modules")
	}

	// Проверяем, что в модулях есть ресурсы
	hasResources := false
	for _, module := range modules {
		resources, err := s.courseRepo.GetResources(ctx, module.ID)
		if err == nil && len(resources) > 0 {
			hasResources = true
			break
		}
	}

	if !hasResources {
		return errors.New("cannot publish course without resources")
	}

	return s.courseRepo.PublishCourse(ctx, courseID, teacherID)
}

func (s *TeacherService) UnpublishCourse(ctx context.Context, courseID, teacherID int64) error {
	return s.courseRepo.UnpublishCourse(ctx, courseID, teacherID)
}

func (s *TeacherService) DeleteCourse(ctx context.Context, courseID, teacherID int64) error {
	// Можно добавить проверку: нельзя удалять опубликованный курс с активными студентами
	return s.courseRepo.Delete(ctx, courseID, teacherID)
}

// ========================================
// СТАТИСТИКА
// ========================================

func (s *TeacherService) GetDashboardStatistics(ctx context.Context, teacherID int64) (*models.TeacherStatistics, error) {
	stats, err := s.statsRepo.GetTeacherStatistics(ctx, teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	// Получаем статистику по курсам
	courseStats, err := s.statsRepo.GetAllCourseStatistics(ctx, teacherID)
	if err == nil {
		stats.CourseStats = courseStats
	}

	return stats, nil
}

func (s *TeacherService) GetCourseStatistics(ctx context.Context, courseID, teacherID int64) (*models.CourseStatistics, error) {
	return s.statsRepo.GetCourseStatistics(ctx, courseID, teacherID)
}

func (s *TeacherService) GetCourseStudentProgress(ctx context.Context, courseID, teacherID int64) ([]*models.StudentCourseProgress, error) {
	return s.statsRepo.GetStudentProgressByCourse(ctx, courseID, teacherID)
}
