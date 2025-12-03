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
	GetLeaderboard(ctx context.Context, limit int) ([]*entities.StudentProfile, error)
	GetLeagueLeaderboard(ctx context.Context, leagueID int, limit int) ([]*entities.StudentProfile, error)
	GetUserGlobalRank(ctx context.Context, userID string) (int, error)
	GetUserLeagueRank(ctx context.Context, userID string) (int, error)
}

type SubjectRepository interface {
	SetInterests(ctx context.Context, userID string, subjectIDs []string) error
	GetByUserID(ctx context.Context, userID string) ([]entities.Subject, error)
}

type ProgressRepository interface {
	GetUserActiveCourses(ctx context.Context, userID string, limit int) ([]entities.CourseProgress, error)
	GetCompletedLessonIDs(ctx context.Context, userID, courseID string) ([]string, error)
	UpsertLessonProgress(ctx context.Context, lp *entities.LessonProgress) error
	GetLessonProgress(ctx context.Context, userID, lessonID string) (*entities.LessonProgress, error)
	CountCompletedLessonsInCourse(ctx context.Context, userID, courseID string) (int, error)
	UpsertCourseProgress(ctx context.Context, cp *entities.CourseProgress) error
	GetAllUserActiveCourses(ctx context.Context, userID string) ([]entities.CourseProgress, error)
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
	GetUserResults(ctx context.Context, userID string) ([]entities.TestResult, error)
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
	AnswerID   string
}

func (s *StudentService) SubmitTest(ctx context.Context, userID, testID string, answers []StudentAnswer) (*entities.TestResult, int, error) {
	test, err := s.testRepo.GetTestFullByID(ctx, testID)
	if err != nil {
		return nil, 0, err
	}

	results, err := s.testRepo.GetUserResults(ctx, userID)
	if err == nil {
		for _, r := range results {
			if r.TestID == testID && r.IsPassed {
				return &r, 0, nil
			}
		}
	}

	correctCount := 0
	totalQuestions := len(test.Questions)

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

	score := 0
	if totalQuestions > 0 {
		score = (correctCount * 100) / totalQuestions
	}
	isPassed := score >= test.PassingScore

	result := &entities.TestResult{
		ID:          uuid.NewString(),
		UserID:      userID,
		TestID:      testID,
		Score:       score,
		IsPassed:    isPassed,
		AttemptDate: time.Now().UTC(),
	}
	if err := s.testRepo.SaveResult(ctx, result); err != nil {
		return nil, 0, err
	}

	xp := 0
	if isPassed {
		xp = 50
		profile, _ := s.profileRepo.GetByUserID(ctx, userID)
		s.addXPToProfile(ctx, profile, int64(xp))
	}

	return result, xp, nil
}

func (s *StudentService) CompleteLesson(ctx context.Context, userID, lessonID string) (*entities.LessonProgress, int, error) {
	lesson, err := s.courseRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, 0, fmt.Errorf("lesson not found: %w", err)
	}

	progress, err := s.progressRepo.GetLessonProgress(ctx, userID, lessonID)
	if err != nil {
		return nil, 0, err
	}

	if progress.IsCompleted {
		return progress, 0, nil
	}

	progress.IsCompleted = true
	progress.Status = entities.StatusCompleted
	progress.LastAccessedAt = time.Now().UTC()

	if err := s.progressRepo.UpsertLessonProgress(ctx, progress); err != nil {
		return nil, 0, err
	}

	xpAwarded := 0
	if lesson.XPReward > 0 {
		profile, err := s.profileRepo.GetByUserID(ctx, userID)
		if err == nil {
			s.updateStreak(profile)

			s.addXPToProfile(ctx, profile, int64(lesson.XPReward))
			xpAwarded = lesson.XPReward
		}
	}

	module, err := s.courseRepo.GetModuleByID(ctx, lesson.ModuleID)
	if err == nil {
		go s.recalculateCourseProgress(context.Background(), userID, module.CourseID)
	}

	return progress, xpAwarded, nil
}

func (s *StudentService) updateStreak(profile *entities.StudentProfile) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if profile.LastActivityDate == nil {
		// Первая активность
		profile.CurrentStreak = 1
		profile.LastActivityDate = &today
	} else {
		lastActivity := time.Date(
			profile.LastActivityDate.Year(),
			profile.LastActivityDate.Month(),
			profile.LastActivityDate.Day(),
			0, 0, 0, 0, time.UTC,
		)

		daysDiff := int(today.Sub(lastActivity).Hours() / 24)

		if daysDiff == 0 {
			// Активность сегодня уже была, ничего не меняем
			return
		} else if daysDiff == 1 {
			// Следующий день подряд - увеличиваем стрик
			profile.CurrentStreak++
			profile.LastActivityDate = &today
		} else {
			// Пропустил дни - стрик сбрасывается
			profile.CurrentStreak = 1
			profile.LastActivityDate = &today
		}
	}

	// Обновляем максимальный стрик
	if profile.CurrentStreak > profile.MaxStreak {
		profile.MaxStreak = profile.CurrentStreak
	}
}

func (s *StudentService) addXPToProfile(ctx context.Context, profile *entities.StudentProfile, xp int64) {
	profile.XP += xp
	profile.WeeklyXP += xp

	// Простой расчет уровня: 100 XP = 1 уровень
	newLevel := int(profile.XP/100) + 1
	if newLevel > profile.Level {
		profile.Level = newLevel
	}

	profile.UpdatedAt = time.Now().UTC()
	s.profileRepo.Update(ctx, profile)
}

func (s *StudentService) recalculateCourseProgress(ctx context.Context, userID, courseID string) {
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

	if err := s.subjectRepo.SetInterests(ctx, userID, subjectIDs); err != nil {
		return fmt.Errorf("failed to set interests: %w", err)
	}

	return nil
}

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
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, nil
	}

	interests, err := s.subjectRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	activeProgress, err := s.progressRepo.GetUserActiveCourses(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	var activeCourses []ActiveCourseData
	for _, prog := range activeProgress {
		course, err := s.courseRepo.GetByID(ctx, prog.CourseID)
		if err == nil {
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

type LeaderboardEntry struct {
	Rank      int
	UserID    string
	FirstName string
	LastName  string
	AvatarURL string
	XP        int64
	Level     int
	LeagueID  int
}

func (s *StudentService) GetWeeklyLeaderboard(ctx context.Context, userID string, limit int) ([]LeaderboardEntry, *int, error) {
	// Получаем лигу пользователя
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Получаем топ игроков в этой лиге по weekly_xp
	profiles, err := s.profileRepo.GetLeagueLeaderboard(ctx, profile.CurrentLeagueID, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get league leaderboard: %w", err)
	}

	// Формируем список с рангами
	entries := make([]LeaderboardEntry, len(profiles))
	userRank := -1

	for i, p := range profiles {
		entries[i] = LeaderboardEntry{
			Rank:     i + 1,
			UserID:   p.UserID,
			XP:       p.WeeklyXP,
			Level:    p.Level,
			LeagueID: p.CurrentLeagueID,
		}

		if p.UserID == userID {
			userRank = i + 1
		}
	}

	// Если пользователь не в топе, получаем его ранг отдельно
	var userRankPtr *int
	if userRank > 0 {
		userRankPtr = &userRank
	} else {
		rank, err := s.profileRepo.GetUserLeagueRank(ctx, userID)
		if err == nil {
			userRankPtr = &rank
		}
	}

	return entries, userRankPtr, nil
}

func (s *StudentService) GetGlobalLeaderboard(ctx context.Context, userID string, limit int) ([]LeaderboardEntry, *int, error) {
	// Получаем топ по общему XP
	profiles, err := s.profileRepo.GetLeaderboard(ctx, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get global leaderboard: %w", err)
	}

	entries := make([]LeaderboardEntry, len(profiles))
	userRank := -1

	for i, p := range profiles {
		entries[i] = LeaderboardEntry{
			Rank:   i + 1,
			UserID: p.UserID,
			XP:     p.XP,
			Level:  p.Level,
		}

		if p.UserID == userID {
			userRank = i + 1
		}
	}

	// Если пользователь не в топе, получаем его ранг отдельно
	var userRankPtr *int
	if userRank > 0 {
		userRankPtr = &userRank
	} else {
		rank, err := s.profileRepo.GetUserGlobalRank(ctx, userID)
		if err == nil {
			userRankPtr = &rank
		}
	}

	return entries, userRankPtr, nil
}

// В интерфейс ProgressRepository добавьте:
// GetAllUserActiveCourses(ctx context.Context, userID string) ([]entities.CourseProgress, error)

// Реализация в StudentService:
func (s *StudentService) GetStudentActiveCourses(ctx context.Context, userID string) ([]ActiveCourseData, error) {
	progressList, err := s.progressRepo.GetAllUserActiveCourses(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []ActiveCourseData
	for _, prog := range progressList {
		course, err := s.courseRepo.GetByID(ctx, prog.CourseID)
		if err == nil {
			result = append(result, ActiveCourseData{
				CourseID:           course.ID,
				Title:              course.Title,
				CoverURL:           course.CoverImageURL,
				ProgressPercentage: prog.ProgressPercentage,
				TotalLessons:       prog.TotalLessonsCount,
				CompletedLessons:   prog.CompletedLessonsCount,
			})
		}
	}
	return result, nil
}
