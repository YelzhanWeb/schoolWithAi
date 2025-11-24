package analytics

import (
	"context"
	"fmt"

	"backend/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsRepository struct {
	connectionURL string
	pool          *pgxpool.Pool
}

func NewAnalyticsRepository(connectionURL string) *AnalyticsRepository {
	return &AnalyticsRepository{connectionURL: connectionURL}
}

func (r *AnalyticsRepository) Connect(ctx context.Context) error {
	p, err := pgxpool.New(ctx, r.connectionURL)
	if err != nil {
		return fmt.Errorf("pgxpool new: %w", err)
	}

	r.pool = p
	return nil
}

func (r *AnalyticsRepository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

// LogActivity сохраняет действие пользователя.
// Пример использования:
// repo.LogActivity(ctx, entities.NewActivityLog(uid, &courseID, "view", nil))
func (r *AnalyticsRepository) LogActivity(ctx context.Context, log *entities.UserActivityLog) error {
	d := newDTO(log)

	query := `
		INSERT INTO user_activity_logs (id, user_id, course_id, action_type, meta_data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		d.ID,
		d.UserID,
		d.CourseID, // pgx сам передаст NULL, если указатель nil
		d.ActionType,
		d.MetaData, // pgx сам сериализует map в jsonb
		d.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}

	return nil
}

// GetUserHistory возвращает последние действия пользователя (например, "Недавно просмотренные").
// Обычно фильтруем только action_type = 'view' или 'complete'.
func (r *AnalyticsRepository) GetUserHistory(
	ctx context.Context,
	userID string,
	limit int,
) ([]*entities.UserActivityLog, error) {
	query := `
		SELECT id, user_id, course_id, action_type, meta_data, created_at
		FROM user_activity_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}
	defer rows.Close()

	var logs []*entities.UserActivityLog

	for rows.Next() {
		var d dto
		// Для сканирования map/jsonb может понадобиться явное приведение,
		// но pgx v5 обычно справляется сам, если структура совпадает.
		err := rows.Scan(
			&d.ID,
			&d.UserID,
			&d.CourseID,
			&d.ActionType,
			&d.MetaData,
			&d.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan log error: %w", err)
		}
		logs = append(logs, d.toEntity())
	}

	return logs, nil
}

// GetUserHistoryByAction позволяет получить конкретные действия (например, только "view")
// Удобно для блока "Вы смотрели"
func (r *AnalyticsRepository) GetUserHistoryByAction(
	ctx context.Context,
	userID, actionType string,
	limit int,
) ([]*entities.UserActivityLog, error) {
	query := `
		SELECT id, user_id, course_id, action_type, meta_data, created_at
		FROM user_activity_logs
		WHERE user_id = $1 AND action_type = $2
		ORDER BY created_at DESC
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, userID, actionType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entities.UserActivityLog
	for rows.Next() {
		var d dto
		rows.Scan(&d.ID, &d.UserID, &d.CourseID, &d.ActionType, &d.MetaData, &d.CreatedAt)
		logs = append(logs, d.toEntity())
	}
	return logs, nil
}
