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

// Start - запускает еженедельный сброс (запускать в main.go как горутину)
func (s *WeeklyResetService) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour) // Проверяем каждый час
	defer ticker.Stop()

	log.Println("Weekly reset scheduler started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Weekly reset scheduler stopped")
			return
		case <-ticker.C:
			now := time.Now().UTC()
			// Проверяем, наступил ли понедельник 00:00 (проверяем диапазон часа)
			if now.Weekday() == time.Monday && now.Hour() == 0 {
				log.Println("Starting weekly reset...")
				s.performWeeklyReset(context.Background())
			}
		}
	}
}

// performWeeklyReset - основная логика сброса
func (s *WeeklyResetService) performWeeklyReset(ctx context.Context) {
	now := time.Now().UTC()
	periodEnd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	periodStart := periodEnd.AddDate(0, 0, -7)

	leagues, err := s.gamificationRepo.GetAllLeagues(ctx)
	if err != nil {
		log.Printf("Failed to get leagues: %v", err)
		return
	}

	// Обрабатываем каждую лигу
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

		// Сохраняем снимки рейтинга
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

		// Продвижение/понижение по лигам
		s.updateLeagues(ctx, profiles, league.ID, len(leagues))
	}

	// Сбрасываем weekly_xp у всех одним запросом
	if err := s.profileRepo.ResetAllWeeklyXP(ctx); err != nil {
		log.Printf("Failed to reset weekly XP: %v", err)
	}

	log.Println("Weekly reset completed successfully!")
}

// updateLeagues - обновление лиг на основе рейтинга
func (s *WeeklyResetService) updateLeagues(ctx context.Context, profiles []*entities.StudentProfile, currentLeagueID, totalLeagues int) {
	if len(profiles) == 0 {
		return
	}

	// ТОП-20% продвигаются вверх (если не максимальная лига)
	promoteCount := len(profiles) / 5
	if promoteCount < 1 {
		promoteCount = 1
	}

	// Нижние 20% опускаются вниз (если не минимальная лига)
	demoteCount := len(profiles) / 5
	if demoteCount < 1 {
		demoteCount = 1
	}

	for i, profile := range profiles {
		shouldUpdate := false

		if i < promoteCount && currentLeagueID < totalLeagues {
			// Продвижение
			profile.CurrentLeagueID++
			shouldUpdate = true
			log.Printf("Promoting user %s to league %d", profile.UserID, profile.CurrentLeagueID)
		} else if i >= len(profiles)-demoteCount && currentLeagueID > 1 {
			// Понижение
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
