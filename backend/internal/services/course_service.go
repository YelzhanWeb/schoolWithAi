package services

import (
	"context"
	"errors"
	"fmt"

	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
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
		fmt.Println(err)

		return nil, errors.New("failed to fetch resources")

	}
	return resources, nil
}

// backend/internal/domain/services/course_service.go
func (s *CourseService) GetResourceByID(ctx context.Context, resourceID int64) (*models.Resource, error) {
	return s.courseRepo.GetResourceByID(ctx, resourceID)
}

func (s *CourseService) GetCoursesByTeacher(ctx context.Context, teacherID int64) ([]*models.Course, error) {
	return s.courseRepo.GetByTeacher(ctx, teacherID)
}

// CreateCourse создает новый курс
func (s *CourseService) CreateCourse(ctx context.Context, course *models.Course) error {
	if err := course.Validate(); err != nil {
		return err
	}
	// Убедимся, что is_published по умолчанию false
	course.IsPublished = false
	return s.courseRepo.Create(ctx, course)
}

// UpdateCourse обновляет курс
func (s *CourseService) UpdateCourse(ctx context.Context, course *models.Course) error {
	// Получаем текущий курс, чтобы не перезаписать IsPublished
	existingCourse, err := s.courseRepo.GetByID(ctx, course.ID)
	if err != nil {
		return errors.New("course not found")
	}

	// Проверяем права
	if existingCourse.CreatedBy != course.CreatedBy {
		return errors.New("permission denied")
	}

	// Обновляем только нужные поля
	existingCourse.Title = course.Title
	existingCourse.Description = course.Description
	existingCourse.DifficultyLevel = course.DifficultyLevel
	existingCourse.AgeGroup = course.AgeGroup
	existingCourse.Subject = course.Subject

	if err := existingCourse.Validate(); err != nil {
		return err
	}

	return s.courseRepo.Update(ctx, existingCourse)
}
