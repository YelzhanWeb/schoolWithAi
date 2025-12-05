package scheduler

import (
	"context"
	"log"
	"time"

	"backend/internal/entities"

	"github.com/google/uuid"
)

type ProfileRepository interface {
	GetLeaderboard(ctx context.Context, limit int) ([]*entities.StudentProfile, error)
	GetLeagueLeaderboard(ctx context.Context, leagueID int, limit int) ([]*entities.StudentProfile, error)
	Update(ctx context.Context, profile *entities.StudentProfile) error
	ResetAllWeeklyXP(ctx context.Context) error
}

type GamificationRepository interface {
	SaveHistorySnapshot(ctx context.Context, h *entities.LeaderboardHistory) error
	GetAllLeagues(ctx context.Context) ([]entities.League, error)
	GetLastResetDate(ctx context.Context) (time.Time, error)
	SetLastResetDate(ctx context.Context, date time.Time) error
}

type WeeklyResetService struct {
	profileRepo      ProfileRepository
	gamificationRepo GamificationRepository
}

func NewWeeklyResetService(pRepo ProfileRepository, gRepo GamificationRepository) *WeeklyResetService {
	return &WeeklyResetService{
		profileRepo:      pRepo,
		gamificationRepo: gRepo,
	}
}

func (s *WeeklyResetService) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Minute) // Проверяем каждые 30 мин
	defer ticker.Stop()

	log.Println("Weekly reset scheduler started")

	// 1. Проверяем сразу при запуске (вдруг сервер лежал во время сброса)
	s.CheckAndRunReset(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 2. И периодически проверяем
			s.CheckAndRunReset(ctx)
		}
	}
}

func (s *WeeklyResetService) CheckAndRunReset(ctx context.Context) {
	now := time.Now().UTC()

	// Вычисляем начало текущей недели (Понедельник 00:00:00)
	// Это наша "Целевая точка сброса"
	currentWeekMonday := getMondayStart(now)

	// Получаем дату последнего успешного сброса из БД
	lastReset, err := s.gamificationRepo.GetLastResetDate(ctx)
	if err != nil {
		log.Printf("Error getting last reset date: %v", err)
		return
	}

	// ЛОГИКА:
	// Если последний сброс был РАНЬШЕ, чем текущий понедельник -> Пора сбрасывать!
	if lastReset.Before(currentWeekMonday) {
		log.Println("New week detected! Starting weekly reset...")

		// Выполняем сброс
		s.performWeeklyReset(ctx)

		// Записываем в БД, что мы выполнили сброс для ЭТОГО понедельника
		if err := s.gamificationRepo.SetLastResetDate(ctx, currentWeekMonday); err != nil {
			log.Printf("Failed to update reset date: %v", err)
		}
	}
}

func getMondayStart(t time.Time) time.Time {
	weekday := t.Weekday()

	daysToSubtract := int(weekday) - int(time.Monday)
	if daysToSubtract < 0 {
		daysToSubtract += 7
	}

	year, month, day := t.AddDate(0, 0, -daysToSubtract).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func (s *WeeklyResetService) performWeeklyReset(ctx context.Context) {
	now := time.Now().UTC()
	periodEnd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	periodStart := periodEnd.AddDate(0, 0, -7)

	leagues, err := s.gamificationRepo.GetAllLeagues(ctx)
	if err != nil {
		log.Printf("Failed to get leagues: %v", err)
		return
	}

	for _, league := range leagues {
		profiles, err := s.profileRepo.GetLeagueLeaderboard(ctx, league.ID, 100)
		if err != nil {
			log.Printf("Failed to get leaderboard for league %d: %v", league.ID, err)
			continue
		}

		if len(profiles) == 0 {
			log.Printf("No profiles in league %d", league.ID)
			continue
		}

		for rank, profile := range profiles {
			history := &entities.LeaderboardHistory{
				ID:          uuid.NewString(),
				PeriodStart: periodStart,
				PeriodEnd:   periodEnd,
				UserID:      profile.UserID,
				LeagueID:    league.ID,
				Rank:        rank + 1,
				TotalXP:     profile.WeeklyXP,
				CreatedAt:   now,
			}
			if err := s.gamificationRepo.SaveHistorySnapshot(ctx, history); err != nil {
				log.Printf("Failed to save history snapshot: %v", err)
			}
		}

		s.updateLeagues(ctx, profiles, league.ID, len(leagues))
	}

	if err := s.profileRepo.ResetAllWeeklyXP(ctx); err != nil {
		log.Printf("Failed to reset weekly XP: %v", err)
	}

	log.Println("Weekly reset completed successfully!")
}

func (s *WeeklyResetService) updateLeagues(ctx context.Context, profiles []*entities.StudentProfile, currentLeagueID, totalLeagues int) {
	if len(profiles) == 0 {
		return
	}

	promoteCount := len(profiles) / 5
	if promoteCount < 1 {
		promoteCount = 1
	}

	demoteCount := len(profiles) / 5
	if demoteCount < 1 {
		demoteCount = 1
	}

	for i, profile := range profiles {
		shouldUpdate := false

		if i < promoteCount && currentLeagueID < totalLeagues {
			profile.CurrentLeagueID++
			shouldUpdate = true
			log.Printf("Promoting user %s to league %d", profile.UserID, profile.CurrentLeagueID)
		} else if i >= len(profiles)-demoteCount && currentLeagueID > 1 {
			profile.CurrentLeagueID--
			shouldUpdate = true
			log.Printf("Demoting user %s to league %d", profile.UserID, profile.CurrentLeagueID)
		}

		if shouldUpdate {
			if err := s.profileRepo.Update(ctx, profile); err != nil {
				log.Printf("Failed to update profile for user %s: %v", profile.UserID, err)
			}
		}
	}
}
