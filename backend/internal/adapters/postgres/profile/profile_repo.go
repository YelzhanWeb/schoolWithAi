package profile

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend/internal/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentProfileRepository struct {
	connectionURL string
	pool          *pgxpool.Pool
}

func NewStudentProfileRepository(connectionURL string) *StudentProfileRepository {
	return &StudentProfileRepository{connectionURL: connectionURL}
}

func (r *StudentProfileRepository) Connect(ctx context.Context) error {
	p, err := pgxpool.New(ctx, r.connectionURL)
	if err != nil {
		return fmt.Errorf("pgxpool new: %w", err)
	}

	r.pool = p
	return nil
}

func (r *StudentProfileRepository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *StudentProfileRepository) Create(ctx context.Context, profile *entities.StudentProfile) error {
	if r.pool == nil {
		return fmt.Errorf("not connected to pool")
	}

	if profile.CreatedAt.IsZero() {
		now := time.Now().UTC()
		profile.CreatedAt = now
		profile.UpdatedAt = now
	}

	d := newDTO(profile)

	query := `
		INSERT INTO student_profiles (
			id, user_id, grade, xp, level, 
			current_league_id, weekly_xp, current_streak, max_streak, last_activity_date,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		d.ID, d.UserID, d.Grade, d.XP, d.Level,
		d.CurrentLeagueID, d.WeeklyXP, d.CurrentStreak, d.MaxStreak, d.LastActivityDate,
		d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entities.ErrAlreadyExists
		}
		return fmt.Errorf("failed to create student profile: %w", err)
	}

	return nil
}

func (r *StudentProfileRepository) GetByUserID(ctx context.Context, userID string) (*entities.StudentProfile, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("not connected to pool")
	}

	query := `
		SELECT sp.id, sp.user_id, sp.grade, sp.xp, sp.level, 
		       sp.current_league_id, sp.weekly_xp, sp.current_streak, sp.max_streak, sp.last_activity_date,
		       sp.created_at, sp.updated_at,
               u.first_name, u.last_name, u.avatar_url
		FROM student_profiles sp
        JOIN users u ON sp.user_id = u.id
		WHERE sp.user_id = $1
	`

	profile, err := scan(r.pool.QueryRow(ctx, query, userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get student profile: %w", err)
	}

	return &profile, nil
}

func (r *StudentProfileRepository) Update(ctx context.Context, profile *entities.StudentProfile) error {
	if r.pool == nil {
		return fmt.Errorf("not connected to pool")
	}

	profile.UpdatedAt = time.Now().UTC()
	d := newDTO(profile)

	query := `
		UPDATE student_profiles
		SET grade = $2,
			xp = $3,
			level = $4,
			current_league_id = $5,
			weekly_xp = $6,
			current_streak = $7,
			max_streak = $8,
			last_activity_date = $9,
			updated_at = $10
		WHERE id = $1
	`

	tag, err := r.pool.Exec(
		ctx,
		query,
		d.ID,
		d.Grade,
		d.XP,
		d.Level,
		d.CurrentLeagueID,
		d.WeeklyXP,
		d.CurrentStreak,
		d.MaxStreak,
		d.LastActivityDate,
		d.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update student profile: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}

	return nil
}

func (r *StudentProfileRepository) Exists(ctx context.Context, userID string) (bool, error) {
	if r.pool == nil {
		return false, fmt.Errorf("not connected to pool")
	}

	query := `SELECT EXISTS(SELECT 1 FROM student_profiles WHERE user_id = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}

	return exists, nil
}

func (r *StudentProfileRepository) GetLeaderboard(ctx context.Context, limit int) ([]*entities.StudentProfile, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("not connected to pool")
	}

	query := `
		SELECT sp.id, sp.user_id, sp.grade, sp.xp, sp.level, 
		       sp.current_league_id, sp.weekly_xp, sp.current_streak, sp.max_streak, sp.last_activity_date,
		       sp.created_at, sp.updated_at,
               u.first_name, u.last_name, u.avatar_url
		FROM student_profiles sp
        JOIN users u ON sp.user_id = u.id
		ORDER BY sp.xp DESC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get global leaderboard: %w", err)
	}
	defer rows.Close()

	var profiles []*entities.StudentProfile

	for rows.Next() {
		profile, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

func (r *StudentProfileRepository) GetLeagueLeaderboard(
	ctx context.Context,
	leagueID int,
	limit int,
) ([]*entities.StudentProfile, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("not connected to pool")
	}

	query := `
		SELECT sp.id, sp.user_id, sp.grade, sp.xp, sp.level, 
		       sp.current_league_id, sp.weekly_xp, sp.current_streak, sp.max_streak, sp.last_activity_date,
		       sp.created_at, sp.updated_at,
               u.first_name, u.last_name, u.avatar_url
		FROM student_profiles sp
        JOIN users u ON sp.user_id = u.id
		WHERE sp.current_league_id = $1
		ORDER BY sp.weekly_xp DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, leagueID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get league leaderboard: %w", err)
	}
	defer rows.Close()

	var profiles []*entities.StudentProfile
	for rows.Next() {
		profile, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

// GetUserGlobalRank возвращает позицию пользователя в глобальном рейтинге
func (r *StudentProfileRepository) GetUserGlobalRank(ctx context.Context, userID string) (int, error) {
	if r.pool == nil {
		return 0, fmt.Errorf("not connected to pool")
	}

	query := `
		WITH ranked_profiles AS (
			SELECT user_id, ROW_NUMBER() OVER (ORDER BY xp DESC) as rank
			FROM student_profiles
		)
		SELECT rank FROM ranked_profiles WHERE user_id = $1
	`

	var rank int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&rank)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, entities.ErrNotFound
		}
		return 0, fmt.Errorf("failed to get user global rank: %w", err)
	}

	return rank, nil
}

// GetUserLeagueRank возвращает позицию пользователя в его лиге
func (r *StudentProfileRepository) GetUserLeagueRank(ctx context.Context, userID string) (int, error) {
	if r.pool == nil {
		return 0, fmt.Errorf("not connected to pool")
	}

	query := `
		WITH user_league AS (
			SELECT current_league_id FROM student_profiles WHERE user_id = $1
		),
		ranked_profiles AS (
			SELECT sp.user_id, ROW_NUMBER() OVER (ORDER BY sp.weekly_xp DESC) as rank
			FROM student_profiles sp
			CROSS JOIN user_league ul
			WHERE sp.current_league_id = ul.current_league_id
		)
		SELECT rank FROM ranked_profiles WHERE user_id = $1
	`

	var rank int
	err := r.pool.QueryRow(ctx, query, userID).Scan(&rank)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, entities.ErrNotFound
		}
		return 0, fmt.Errorf("failed to get user league rank: %w", err)
	}

	return rank, nil
}

// ResetAllWeeklyXP сбрасывает weekly_xp у всех профилей (для еженедельного сброса)
func (r *StudentProfileRepository) ResetAllWeeklyXP(ctx context.Context) error {
	if r.pool == nil {
		return fmt.Errorf("not connected to pool")
	}

	query := `UPDATE student_profiles SET weekly_xp = 0, updated_at = NOW()`

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to reset weekly XP: %w", err)
	}

	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scan(scanner rowScanner) (entities.StudentProfile, error) {
	var d dto
	err := scanner.Scan(
		&d.ID,
		&d.UserID,
		&d.Grade,
		&d.XP,
		&d.Level,
		&d.CurrentLeagueID,
		&d.WeeklyXP,
		&d.CurrentStreak,
		&d.MaxStreak,
		&d.LastActivityDate,
		&d.CreatedAt,
		&d.UpdatedAt,
		&d.FirstName,
		&d.LastName,
		&d.AvatarURL,
	)
	if err != nil {
		return entities.StudentProfile{}, err
	}

	return *d.toEntity(), nil
}
