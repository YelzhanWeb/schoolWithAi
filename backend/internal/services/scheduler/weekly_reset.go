// backend/internal/services/scheduler/weekly_reset.go
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
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().UTC()
			// Проверяем, наступил ли понедельник 00:00
			if now.Weekday() == time.Monday && now.Hour() == 0 {
				s.performWeeklyReset(context.Background())
			}
		}
	}
}

// performWeeklyReset - основная логика сброса
func (s *WeeklyResetService) performWeeklyReset(ctx context.Context) {
	log.Println("Starting weekly reset...")

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
			s.gamificationRepo.SaveHistorySnapshot(ctx, history)
		}

		// Продвижение/понижение по лигам
		s.updateLeagues(ctx, profiles, league.ID, len(leagues))
	}

	// Сбрасываем weekly_xp у всех
	s.resetWeeklyXP(ctx)

	log.Println("Weekly reset completed!")
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
		if i < promoteCount && currentLeagueID < totalLeagues {
			// Продвижение
			profile.CurrentLeagueID++
			s.profileRepo.Update(ctx, profile)
		} else if i >= len(profiles)-demoteCount && currentLeagueID > 1 {
			// Понижение
			profile.CurrentLeagueID--
			s.profileRepo.Update(ctx, profile)
		}
	}
}

// resetWeeklyXP - сбрасываем weekly_xp у всех студентов
func (s *WeeklyResetService) resetWeeklyXP(ctx context.Context) {
	// Получаем всех (можно оптимизировать через батч-запрос в БД)
	profiles, err := s.profileRepo.GetLeaderboard(ctx, 10000)
	if err != nil {
		log.Printf("Failed to get all profiles: %v", err)
		return
	}

	for _, profile := range profiles {
		profile.WeeklyXP = 0
		s.profileRepo.Update(ctx, profile)
	}
}
