package entities

import (
	"time"
)

type League struct {
	ID         int
	Slug       string
	Name       string
	OrderIndex int
	IconURL    string
}

type Achievement struct {
	ID          string
	Slug        string
	Name        string
	Description string
	IconURL     string
	XPReward    int
}

type UserAchievement struct {
	UserID        string
	AchievementID string
	EarnedAt      time.Time

	Achievement *Achievement
}

type LeaderboardHistory struct {
	ID          string
	PeriodStart time.Time
	PeriodEnd   time.Time
	UserID      string
	LeagueID    int
	Rank        int
	TotalXP     int64
	CreatedAt   time.Time
}
