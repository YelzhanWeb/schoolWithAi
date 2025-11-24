package gamification

import (
	"context"
	"fmt"
	"time"

	"backend/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GamificationRepository struct {
	connectionURL string
	pool          *pgxpool.Pool
}

func NewGamificationRepository(connectionURL string) *GamificationRepository {
	return &GamificationRepository{connectionURL: connectionURL}
}

func (r *GamificationRepository) Connect(ctx context.Context) error {
	p, err := pgxpool.New(ctx, r.connectionURL)
	if err != nil {
		return fmt.Errorf("pgxpool new: %w", err)
	}

	r.pool = p
	return nil
}

func (r *GamificationRepository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *GamificationRepository) GetAllLeagues(ctx context.Context) ([]entities.League, error) {
	query := `SELECT id, slug, name, order_index, icon_url FROM leagues ORDER BY order_index ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get leagues: %w", err)
	}
	defer rows.Close()

	var list []entities.League
	for rows.Next() {
		var d leagueDTO
		if err := rows.Scan(&d.ID, &d.Slug, &d.Name, &d.OrderIndex, &d.IconURL); err != nil {
			return nil, err
		}
		list = append(list, d.toEntity())
	}
	return list, nil
}

func (r *GamificationRepository) GetAllAchievements(ctx context.Context) ([]entities.Achievement, error) {
	query := `SELECT id, slug, name, description, icon_url, xp_reward FROM achievements`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all achievements: %w", err)
	}
	defer rows.Close()

	var list []entities.Achievement
	for rows.Next() {
		var d achievementDTO
		if err := rows.Scan(&d.ID, &d.Slug, &d.Name, &d.Description, &d.IconURL, &d.XPReward); err != nil {
			return nil, err
		}
		list = append(list, d.toEntity())
	}
	return list, nil
}

func (r *GamificationRepository) GetUserAchievements(
	ctx context.Context,
	userID string,
) ([]entities.UserAchievement, error) {
	query := `
		SELECT ua.user_id, ua.achievement_id, ua.earned_at,
		       a.slug, a.name, a.description, a.icon_url, a.xp_reward
		FROM user_achievements ua
		JOIN achievements a ON ua.achievement_id = a.id
		WHERE ua.user_id = $1
		ORDER BY ua.earned_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get user achievements: %w", err)
	}
	defer rows.Close()

	var list []entities.UserAchievement
	for rows.Next() {
		var d userAchievementDTO
		err := rows.Scan(
			&d.UserID, &d.AchievementID, &d.EarnedAt,
			&d.Slug, &d.Name, &d.Description, &d.IconURL, &d.XPReward,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, d.toEntity())
	}
	return list, nil
}

func (r *GamificationRepository) GiveAchievement(ctx context.Context, userID, achievementSlug string) (bool, error) {
	var achievementID string
	err := r.pool.QueryRow(ctx, `SELECT id FROM achievements WHERE slug = $1`, achievementSlug).Scan(&achievementID)
	if err != nil {
		return false, fmt.Errorf("achievement slug '%s' not found: %w", achievementSlug, err)
	}

	query := `
		INSERT INTO user_achievements (user_id, achievement_id, earned_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id, achievement_id) DO NOTHING
	`

	tag, err := r.pool.Exec(ctx, query, userID, achievementID)
	if err != nil {
		return false, fmt.Errorf("failed to give achievement: %w", err)
	}

	return tag.RowsAffected() > 0, nil
}

func (r *GamificationRepository) SaveHistorySnapshot(ctx context.Context, h *entities.LeaderboardHistory) error {
	query := `
		INSERT INTO leaderboard_history (id, period_start, period_end, user_id, league_id, rank, total_xp, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(
		ctx, query,
		h.ID, h.PeriodStart, h.PeriodEnd, h.UserID, h.LeagueID, h.Rank, h.TotalXP, h.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("save history: %w", err)
	}
	return nil
}

func (r *GamificationRepository) GetHistoryByPeriod(
	ctx context.Context,
	start, end time.Time,
	leagueID int,
) ([]entities.LeaderboardHistory, error) {
	query := `
		SELECT id, period_start, period_end, user_id, league_id, rank, total_xp, created_at
		FROM leaderboard_history
		WHERE period_start = $1 AND period_end = $2 AND league_id = $3
		ORDER BY rank ASC
	`

	rows, err := r.pool.Query(ctx, query, start, end, leagueID)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}
	defer rows.Close()

	var list []entities.LeaderboardHistory
	for rows.Next() {
		var d historyDTO
		err := rows.Scan(&d.ID, &d.PeriodStart, &d.PeriodEnd, &d.UserID, &d.LeagueID, &d.Rank, &d.TotalXP, &d.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, d.toEntity())
	}
	return list, nil
}
