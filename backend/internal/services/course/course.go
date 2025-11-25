package course

import (
	"context"
	"fmt"

	"backend/internal/entities"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entities.Course) error
	UpdateCourse(ctx context.Context, course *entities.Course) error
	GetByID(ctx context.Context, id string) (*entities.Course, error)
}

type CourseService struct {
	repo CourseRepository
}

func NewCourseService(repo CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) CreateCourse(ctx context.Context, course *entities.Course) error {
	if err := s.repo.Create(ctx, course); err != nil {
		return fmt.Errorf("failed to create course: %w", err)
	}
	return nil
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
