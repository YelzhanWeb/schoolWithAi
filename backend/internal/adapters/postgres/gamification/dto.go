package gamification

import (
	"time"

	"backend/internal/entities"
)

type leagueDTO struct {
	ID         int    `db:"id"`
	Slug       string `db:"slug"`
	Name       string `db:"name"`
	OrderIndex int    `db:"order_index"`
	IconURL    string `db:"icon_url"`
}

func (d *leagueDTO) toEntity() entities.League {
	return entities.League{
		ID:         d.ID,
		Slug:       d.Slug,
		Name:       d.Name,
		OrderIndex: d.OrderIndex,
		IconURL:    d.IconURL,
	}
}

type achievementDTO struct {
	ID          string `db:"id"`
	Slug        string `db:"slug"`
	Name        string `db:"name"`
	Description string `db:"description"`
	IconURL     string `db:"icon_url"`
	XPReward    int    `db:"xp_reward"`
}

func (d *achievementDTO) toEntity() entities.Achievement {
	return entities.Achievement{
		ID:          d.ID,
		Slug:        d.Slug,
		Name:        d.Name,
		Description: d.Description,
		IconURL:     d.IconURL,
		XPReward:    d.XPReward,
	}
}

type userAchievementDTO struct {
	UserID        string    `db:"user_id"`
	AchievementID string    `db:"achievement_id"`
	EarnedAt      time.Time `db:"earned_at"`

	Slug        string `db:"slug"`
	Name        string `db:"name"`
	Description string `db:"description"`
	IconURL     string `db:"icon_url"`
	XPReward    int    `db:"xp_reward"`
}

func (d *userAchievementDTO) toEntity() entities.UserAchievement {
	return entities.UserAchievement{
		UserID:        d.UserID,
		AchievementID: d.AchievementID,
		EarnedAt:      d.EarnedAt.UTC(),
		Achievement: &entities.Achievement{
			ID:          d.AchievementID,
			Slug:        d.Slug,
			Name:        d.Name,
			Description: d.Description,
			IconURL:     d.IconURL,
			XPReward:    d.XPReward,
		},
	}
}

type historyDTO struct {
	ID          string    `db:"id"`
	PeriodStart time.Time `db:"period_start"`
	PeriodEnd   time.Time `db:"period_end"`
	UserID      string    `db:"user_id"`
	LeagueID    int       `db:"league_id"`
	Rank        int       `db:"rank"`
	TotalXP     int64     `db:"total_xp"`
	CreatedAt   time.Time `db:"created_at"`
}

func (d *historyDTO) toEntity() entities.LeaderboardHistory {
	return entities.LeaderboardHistory{
		ID:          d.ID,
		PeriodStart: d.PeriodStart,
		PeriodEnd:   d.PeriodEnd,
		UserID:      d.UserID,
		LeagueID:    d.LeagueID,
		Rank:        d.Rank,
		TotalXP:     d.TotalXP,
		CreatedAt:   d.CreatedAt.UTC(),
	}
}
