package course

import (
	"context"
	"fmt"

	"backend/internal/entities"

	"github.com/rs/zerolog/log"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entities.Course) error
	UpdateCourse(ctx context.Context, course *entities.Course) error
	GetByID(ctx context.Context, id string) (*entities.Course, error)
	GetCourseStructure(ctx context.Context, courseID string) ([]entities.Module, error)
	GetByAuthorID(ctx context.Context, authorID string) ([]entities.Course, error)
	DeleteCourse(ctx context.Context, id string) error
	GetCatalog(ctx context.Context) ([]entities.Course, error)

	AddModule(ctx context.Context, module *entities.Module) error
	GetModuleByID(ctx context.Context, moduleID string) (*entities.Module, error) // <-- Добавили
	UpdateModule(ctx context.Context, module *entities.Module) error
	DeleteModule(ctx context.Context, id string) error

	AddLesson(ctx context.Context, lesson *entities.Lesson) error
	GetLessonByID(ctx context.Context, lessonID string) (*entities.Lesson, error)
	UpdateLesson(ctx context.Context, lesson *entities.Lesson) error
	DeleteLesson(ctx context.Context, id string) error

	GetAllTags(ctx context.Context) ([]entities.Tag, error)

	ToggleFavorite(ctx context.Context, userID, courseID string) (bool, error)
	GetUserFavorites(ctx context.Context, userID string) ([]entities.Course, error)
	IsFavorite(ctx context.Context, userID, courseID string) (bool, error)

	GetCoursesByIDs(ctx context.Context, ids []string) ([]entities.Course, error)
}

type MLClient interface {
	GetRecommendedCourseIDs(userID string) ([]string, error)
}

type CourseService struct {
	repo     CourseRepository
	mlClient MLClient
}

func NewCourseService(repo CourseRepository, mlClient MLClient) *CourseService {
	return &CourseService{
		repo:     repo,
		mlClient: mlClient,
	}
}

func (s *CourseService) CreateCourse(ctx context.Context, course *entities.Course) error {
	if err := s.repo.Create(ctx, course); err != nil {
		return fmt.Errorf("failed to create course: %w", err)
	}
	return nil
}

func (s *CourseService) GetCatalog(ctx context.Context) ([]entities.Course, error) {
	return s.repo.GetCatalog(ctx)
}

func (s *CourseService) GetCoursesByAuthor(ctx context.Context, authorID string) ([]entities.Course, error) {
	return s.repo.GetByAuthorID(ctx, authorID)
}

func (s *CourseService) GetCourseByID(ctx context.Context, courseID string) (*entities.Course, error) {
	return s.repo.GetByID(ctx, courseID)
}

func (s *CourseService) GetUserCourse(ctx context.Context, courseID, userID string) (*entities.Course, error) {
	course, err := s.repo.GetByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	if course.AuthorID != userID {
		return nil, fmt.Errorf("access denied: user is not the author")
	}

	return course, nil
}

func (s *CourseService) UpdateCourse(ctx context.Context, courseID string, updates *entities.Course) error {
	existing, err := s.repo.GetByID(ctx, courseID)
	if err != nil {
		return err
	}

	existing.Title = updates.Title
	existing.Description = updates.Description
	existing.DifficultyLevel = updates.DifficultyLevel
	existing.CoverImageURL = updates.CoverImageURL
	existing.SubjectID = updates.SubjectID

	return s.repo.UpdateCourse(ctx, existing)
}

func (s *CourseService) ChangePublishStatus(ctx context.Context, courseID string, isPublished bool) error {
	course, err := s.repo.GetByID(ctx, courseID)
	if err != nil {
		return err
	}

	course.IsPublished = isPublished
	return s.repo.UpdateCourse(ctx, course)
}

func (s *CourseService) DeleteCourse(ctx context.Context, id string) error {
	return s.repo.DeleteCourse(ctx, id)
}

func (s *CourseService) CreateModule(ctx context.Context, userID string, module *entities.Module) error {
	if _, err := s.GetUserCourse(ctx, module.CourseID, userID); err != nil {
		return err
	}
	return s.repo.AddModule(ctx, module)
}

func (s *CourseService) UpdateModule(ctx context.Context, userID string, module *entities.Module) error {
	existing, err := s.repo.GetModuleByID(ctx, module.ID)
	if err != nil {
		return err
	}
	if _, err := s.GetUserCourse(ctx, existing.CourseID, userID); err != nil {
		return err
	}

	existing.Title = module.Title
	existing.OrderIndex = module.OrderIndex

	return s.repo.UpdateModule(ctx, existing)
}

func (s *CourseService) DeleteModule(ctx context.Context, userID, moduleID string) error {
	existing, err := s.repo.GetModuleByID(ctx, moduleID)
	if err != nil {
		return err
	}
	if _, err := s.GetUserCourse(ctx, existing.CourseID, userID); err != nil {
		return err
	}
	return s.repo.DeleteModule(ctx, moduleID)
}

func (s *CourseService) CreateLesson(ctx context.Context, userID string, lesson *entities.Lesson) error {
	module, err := s.repo.GetModuleByID(ctx, lesson.ModuleID)
	if err != nil {
		return fmt.Errorf("module not found: %w", err)
	}

	if _, err := s.GetUserCourse(ctx, module.CourseID, userID); err != nil {
		return err
	}

	return s.repo.AddLesson(ctx, lesson)
}

func (s *CourseService) GetLessonByID(ctx context.Context, lessonID string) (*entities.Lesson, error) {
	return s.repo.GetLessonByID(ctx, lessonID)
}

func (s *CourseService) UpdateLesson(ctx context.Context, userID string, lesson *entities.Lesson) error {
	existing, err := s.repo.GetLessonByID(ctx, lesson.ID)
	if err != nil {
		return err
	}

	module, err := s.repo.GetModuleByID(ctx, existing.ModuleID)
	if err != nil {
		return err
	}
	if _, err := s.GetUserCourse(ctx, module.CourseID, userID); err != nil {
		return err
	}

	existing.Title = lesson.Title
	existing.ContentText = lesson.ContentText
	existing.VideoURL = lesson.VideoURL
	existing.FileAttachmentURL = lesson.FileAttachmentURL
	existing.OrderIndex = lesson.OrderIndex
	existing.XPReward = lesson.XPReward

	return s.repo.UpdateLesson(ctx, existing)
}

func (s *CourseService) DeleteLesson(ctx context.Context, userID, lessonID string) error {
	existing, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return err
	}
	module, err := s.repo.GetModuleByID(ctx, existing.ModuleID)
	if err != nil {
		return err
	}
	if _, err := s.GetUserCourse(ctx, module.CourseID, userID); err != nil {
		return err
	}

	return s.repo.DeleteLesson(ctx, lessonID)
}

func (s *CourseService) GetFullStructure(ctx context.Context, courseID string) ([]entities.Module, error) {
	return s.repo.GetCourseStructure(ctx, courseID)
}

func (s *CourseService) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	return s.repo.GetAllTags(ctx)
}

func (s *CourseService) ToggleFavorite(ctx context.Context, userID, courseID string) (bool, error) {
	return s.repo.ToggleFavorite(ctx, userID, courseID)
}

func (s *CourseService) GetUserFavorites(ctx context.Context, userID string) ([]entities.Course, error) {
	return s.repo.GetUserFavorites(ctx, userID)
}

func (s *CourseService) IsCourseFavorite(ctx context.Context, userID, courseID string) (bool, error) {
	return s.repo.IsFavorite(ctx, userID, courseID)
}

func (s *CourseService) GetRecommendations(ctx context.Context, userID string) ([]entities.Course, error) {
	courseIDs, err := s.mlClient.GetRecommendedCourseIDs(userID)
	if err != nil {
		log.Warn().Err(err).Msg("Warning: ML service failed")
		return []entities.Course{}, nil
	}

	if len(courseIDs) == 0 {
		return []entities.Course{}, nil
	}

	courses, err := s.repo.GetCoursesByIDs(ctx, courseIDs)
	if err != nil {
		return nil, err
	}

	// Опционально: можно восстановить порядок курсов, как вернул ML (база может вернуть вразнобой)
	// Но для начала и так сойдет.

	return courses, nil
}
