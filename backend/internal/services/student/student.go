package student

import (
	"context"
	"fmt"

	"backend/internal/entities"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *entities.StudentProfile) error
	GetByUserID(ctx context.Context, userID string) (*entities.StudentProfile, error)
	Update(ctx context.Context, profile *entities.StudentProfile) error
	Exists(ctx context.Context, userID string) (bool, error)
}

type SubjectRepository interface {
	SetInterests(ctx context.Context, userID string, subjectIDs []string) error
	GetByUserID(ctx context.Context, userID string) ([]entities.Subject, error)
}

type ProgressRepository interface {
	GetUserActiveCourses(ctx context.Context, userID string, limit int) ([]entities.CourseProgress, error)
}

type CourseRepository interface {
	GetByID(ctx context.Context, id string) (*entities.Course, error)
}

type GamificationRepository interface {
	GetAllLeagues(ctx context.Context) ([]entities.League, error)
}

type StudentService struct {
	profileRepo      ProfileRepository
	subjectRepo      SubjectRepository
	progressRepo     ProgressRepository
	courseRepo       CourseRepository
	gamificationRepo GamificationRepository
}

func NewStudentService(
	pRepo ProfileRepository,
	sRepo SubjectRepository,
	prRepo ProgressRepository,
	cRepo CourseRepository,
	gRepo GamificationRepository,
) *StudentService {
	return &StudentService{
		profileRepo:      pRepo,
		subjectRepo:      sRepo,
		progressRepo:     prRepo,
		courseRepo:       cRepo,
		gamificationRepo: gRepo,
	}
}

// --- ONBOARDING ---

func (s *StudentService) CompleteOnboarding(ctx context.Context, userID string, grade int, subjectIDs []string) error {
	// 1. Создаем или обновляем профиль (класс)
	exists, err := s.profileRepo.Exists(ctx, userID)
	if err != nil {
		return err
	}

	if exists {
		profile, err := s.profileRepo.GetByUserID(ctx, userID)
		if err != nil {
			return err
		}
		profile.Grade = grade
		if err := s.profileRepo.Update(ctx, profile); err != nil {
			return err
		}
	} else {
		profile, err := entities.NewStudentProfile(userID, grade)
		if err != nil {
			return err
		}
		if err := s.profileRepo.Create(ctx, profile); err != nil {
			return err
		}
	}

	// 2. Сохраняем интересы
	if err := s.subjectRepo.SetInterests(ctx, userID, subjectIDs); err != nil {
		return fmt.Errorf("failed to set interests: %w", err)
	}

	return nil
}

// --- DASHBOARD DATA ---

type DashboardData struct {
	Profile       *entities.StudentProfile
	Interests     []entities.Subject
	ActiveCourses []ActiveCourseData
}

type ActiveCourseData struct {
	CourseID           string
	Title              string
	CoverURL           string
	ProgressPercentage int
	TotalLessons       int
	CompletedLessons   int
}

func (s *StudentService) GetDashboardData(ctx context.Context, userID string) (*DashboardData, error) {
	// 1. Профиль (XP, уровень)
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Если профиля нет, возвращаем nil (фронт должен отправить на онбординг)
		return nil, nil
	}

	// 2. Интересы
	interests, err := s.subjectRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. Активные курсы (с прогрессом)
	activeProgress, err := s.progressRepo.GetUserActiveCourses(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	// Обогащаем прогресс данными о курсе (название, картинка)
	var activeCourses []ActiveCourseData
	for _, prog := range activeProgress {
		course, err := s.courseRepo.GetByID(ctx, prog.CourseID)
		if err == nil { // Если курс не найден, просто пропускаем
			activeCourses = append(activeCourses, ActiveCourseData{
				CourseID:           course.ID,
				Title:              course.Title,
				CoverURL:           course.CoverImageURL,
				ProgressPercentage: prog.ProgressPercentage,
				TotalLessons:       prog.TotalLessonsCount,
				CompletedLessons:   prog.CompletedLessonsCount,
			})
		}
	}

	return &DashboardData{
		Profile:       profile,
		Interests:     interests,
		ActiveCourses: activeCourses,
	}, nil
}
