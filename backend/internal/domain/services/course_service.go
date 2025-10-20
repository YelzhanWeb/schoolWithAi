package services

import (
	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
	"context"
	"errors"
)

type CourseService struct {
	courseRepo repositories.CourseRepository
}

func NewCourseService(courseRepo repositories.CourseRepository) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
	}
}

// GetAllCourses - бизнес-логика получения всех курсов
func (s *CourseService) GetAllCourses(ctx context.Context) ([]*models.Course, error) {
	// Здесь может быть дополнительная бизнес-логика:
	// - Фильтрация по правам
	// - Кеширование
	// - Логирование

	courses, err := s.courseRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch courses")
	}

	return courses, nil
}

// GetCourseDetails - получить курс со всеми модулями
func (s *CourseService) GetCourseDetails(ctx context.Context, courseID int64) (*models.Course, []*models.Module, error) {
	// Получить курс
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, nil, errors.New("course not found")
	}

	// Получить модули
	modules, err := s.courseRepo.GetModules(ctx, courseID)
	if err != nil {
		return nil, nil, errors.New("failed to fetch modules")
	}

	return course, modules, nil
}

// GetModuleResources - получить ресурсы модуля
func (s *CourseService) GetModuleResources(ctx context.Context, moduleID int64) ([]*models.Resource, error) {
	resources, err := s.courseRepo.GetResources(ctx, moduleID)
	if err != nil {
		return nil, errors.New("failed to fetch resources")
	}

	return resources, nil
}

// backend/internal/domain/services/course_service.go
func (s *CourseService) GetResourceByID(ctx context.Context, resourceID int64) (*models.Resource, error) {
	return s.courseRepo.GetResourceByID(ctx, resourceID)
}
