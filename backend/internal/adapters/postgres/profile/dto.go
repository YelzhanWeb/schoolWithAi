package profile

import (
	"time"

	"backend/internal/entities"
)

type dto struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`
	Grade  int    `db:"grade"`
	XP     int64  `db:"xp"`
	Level  int    `db:"level"`

	CurrentLeagueID  int        `db:"current_league_id"`
	WeeklyXP         int64      `db:"weekly_xp"`
	CurrentStreak    int        `db:"current_streak"`
	MaxStreak        int        `db:"max_streak"`
	LastActivityDate *time.Time `db:"last_activity_date"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func newDTO(sp *entities.StudentProfile) dto {
	return dto{
		ID:               sp.ID,
		UserID:           sp.UserID,
		Grade:            sp.Grade,
		XP:               sp.XP,
		Level:            sp.Level,
		CurrentLeagueID:  sp.CurrentLeagueID,
		WeeklyXP:         sp.WeeklyXP,
		CurrentStreak:    sp.CurrentStreak,
		MaxStreak:        sp.MaxStreak,
		LastActivityDate: sp.LastActivityDate,
		CreatedAt:        sp.CreatedAt,
		UpdatedAt:        sp.UpdatedAt,
	}
}

func (d *dto) toEntity() *entities.StudentProfile {
	return &entities.StudentProfile{
		ID:               d.ID,
		UserID:           d.UserID,
		Grade:            d.Grade,
		XP:               d.XP,
		Level:            d.Level,
		CurrentLeagueID:  d.CurrentLeagueID,
		WeeklyXP:         d.WeeklyXP,
		CurrentStreak:    d.CurrentStreak,
		MaxStreak:        d.MaxStreak,
		LastActivityDate: d.LastActivityDate,
		CreatedAt:        d.CreatedAt.UTC(),
		UpdatedAt:        d.UpdatedAt.UTC(),
	}
}
