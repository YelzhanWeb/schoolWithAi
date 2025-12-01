package student

import (
	"context"
	"fmt"
	"time"

	"backend/internal/entities"

	"github.com/google/uuid"
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
	GetCompletedLessonIDs(ctx context.Context, userID, courseID string) ([]string, error)
	UpsertLessonProgress(ctx context.Context, lp *entities.LessonProgress) error
	GetLessonProgress(
		ctx context.Context,
		userID, lessonID string,
	) (*entities.LessonProgress, error)
	CountCompletedLessonsInCourse(ctx context.Context, userID, courseID string) (int, error)
	UpsertCourseProgress(ctx context.Context, cp *entities.CourseProgress) error
}

type CourseRepository interface {
	GetByID(ctx context.Context, id string) (*entities.Course, error)
	GetLessonByID(ctx context.Context, id string) (*entities.Lesson, error)
	GetModuleByID(ctx context.Context, id string) (*entities.Module, error)
	GetCourseStructure(ctx context.Context, courseID string) ([]entities.Module, error)
}

type GamificationRepository interface {
	GetAllLeagues(ctx context.Context) ([]entities.League, error)
}

type TestRepository interface {
	GetTestByModuleID(ctx context.Context, moduleID string) (*entities.Test, error)
	SaveResult(ctx context.Context, res *entities.TestResult) error
	GetTestFullByID(ctx context.Context, testID string) (*entities.Test, error)
}

type StudentService struct {
	profileRepo      ProfileRepository
	subjectRepo      SubjectRepository
	progressRepo     ProgressRepository
	courseRepo       CourseRepository
	gamificationRepo GamificationRepository
	testRepo         TestRepository
}

func NewStudentService(
	pRepo ProfileRepository,
	sRepo SubjectRepository,
	prRepo ProgressRepository,
	cRepo CourseRepository,
	gRepo GamificationRepository,
	tRepo TestRepository,
) *StudentService {
	return &StudentService{
		profileRepo:      pRepo,
		subjectRepo:      sRepo,
		progressRepo:     prRepo,
		courseRepo:       cRepo,
		gamificationRepo: gRepo,
		testRepo:         tRepo,
	}
}

type StudentAnswer struct {
	QuestionID string
	AnswerID   string // Если single choice
	// AnswerIDs []string // Если multiple choice (пока пропустим для простоты)
}

func (s *StudentService) SubmitTest(ctx context.Context, userID, testID string, answers []StudentAnswer) (*entities.TestResult, int, error) {
	// 1. Нам нужно получить "Правильные ответы" из БД.
	// В текущем TestRepo (testing.go) нет метода GetTestByID.
	// ДАВАЙ ИСПОЛЬЗОВАТЬ GetTestByModuleID, но нам нужно знать moduleID.
	// Это неудобно. Давай добавим GetTestByIDWithAnswers в TestRepo.

	// ПРЕДПОЛОЖИМ, что мы его добавили (см. шаг ниже)
	test, err := s.testRepo.GetTestFullByID(ctx, testID)
	if err != nil {
		return nil, 0, err
	}

	correctCount := 0
	totalQuestions := len(test.Questions)

	// 2. Проверяем ответы
	// Создаем map правильных ответов для быстрого поиска: QuestionID -> CorrectAnswerID
	correctMap := make(map[string]string)
	for _, q := range test.Questions {
		for _, a := range q.Answers {
			if a.IsCorrect {
				correctMap[q.ID] = a.ID
				break
			}
		}
	}

	for _, userAns := range answers {
		if correctID, ok := correctMap[userAns.QuestionID]; ok {
			if correctID == userAns.AnswerID {
				correctCount++
			}
		}
	}

	// 3. Считаем результат
	score := 0
	if totalQuestions > 0 {
		score = (correctCount * 100) / totalQuestions
	}
	isPassed := score >= test.PassingScore

	// 4. Сохраняем результат
	result := &entities.TestResult{
		ID:          uuid.NewString(), // Не забудь import "github.com/google/uuid"
		UserID:      userID,
		TestID:      testID,
		Score:       score,
		IsPassed:    isPassed,
		AttemptDate: time.Now().UTC(),
	}
	if err := s.testRepo.SaveResult(ctx, result); err != nil {
		return nil, 0, err
	}

	// 5. Начисляем XP если сдал (например, 50 XP за тест)
	xp := 0
	if isPassed {
		// Проверяем, сдавал ли он раньше этот тест (чтобы не фармить XP)
		// (Это сложнее, пока пропустим, дадим XP за каждую успешную сдачу для радости)
		xp = 50
		profile, _ := s.profileRepo.GetByUserID(ctx, userID)
		profile.AddXP(int64(xp))
		s.profileRepo.Update(ctx, profile)
	}

	return result, xp, nil
}

func (s *StudentService) CompleteLesson(ctx context.Context, userID, lessonID string) (*entities.LessonProgress, int, error) {
	// А. Получаем данные урока
	lesson, err := s.courseRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, 0, fmt.Errorf("lesson not found: %w", err)
	}

	// Б. Проверяем, не пройден ли он уже
	progress, err := s.progressRepo.GetLessonProgress(ctx, userID, lessonID)
	if err != nil {
		return nil, 0, err
	}

	if progress.IsCompleted {
		// Если уже пройден, просто возвращаем (XP повторно не даем)
		return progress, 0, nil
	}

	// В. Обновляем статус
	progress.IsCompleted = true
	progress.Status = entities.StatusCompleted
	progress.LastAccessedAt = time.Now().UTC()

	if err := s.progressRepo.UpsertLessonProgress(ctx, progress); err != nil {
		return nil, 0, err
	}

	// Г. Начисляем XP (если есть награда)
	xpAwarded := 0
	if lesson.XPReward > 0 {
		profile, err := s.profileRepo.GetByUserID(ctx, userID)
		if err == nil {
			profile.AddXP(int64(lesson.XPReward))
			profile.WeeklyXP += int64(lesson.XPReward) // Обновляем и недельный рейтинг
			s.profileRepo.Update(ctx, profile)
			xpAwarded = lesson.XPReward
		}
	}

	// Д. Пересчитываем общий прогресс курса
	// Нам нужно найти course_id через модуль
	module, err := s.courseRepo.GetModuleByID(ctx, lesson.ModuleID)
	if err == nil {
		go s.recalculateCourseProgress(context.Background(), userID, module.CourseID)
	}

	return progress, xpAwarded, nil
}

// Вспомогательная функция для пересчета %
func (s *StudentService) recalculateCourseProgress(ctx context.Context, userID, courseID string) {
	// 1. Считаем всего уроков
	modules, err := s.courseRepo.GetCourseStructure(ctx, courseID)
	if err != nil {
		return
	}

	totalLessons := 0
	for _, m := range modules {
		totalLessons += len(m.Lessons)
	}
	if totalLessons == 0 {
		return
	}

	// 2. Считаем пройденные
	completedCount, err := s.progressRepo.CountCompletedLessonsInCourse(ctx, userID, courseID)
	if err != nil {
		return
	}

	// 3. Вычисляем процент
	percentage := (completedCount * 100) / totalLessons
	isCompleted := percentage == 100

	// 4. Сохраняем
	progress := &entities.CourseProgress{
		UserID:                userID,
		CourseID:              courseID,
		CompletedLessonsCount: completedCount,
		TotalLessonsCount:     totalLessons,
		ProgressPercentage:    percentage,
		IsCompleted:           isCompleted,
		UpdatedAt:             time.Now().UTC(),
	}
	s.progressRepo.UpsertCourseProgress(ctx, progress)
}

func (s *StudentService) GetCourseProgress(ctx context.Context, userID, courseID string) ([]string, error) {
	return s.progressRepo.GetCompletedLessonIDs(ctx, userID, courseID)
}

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
