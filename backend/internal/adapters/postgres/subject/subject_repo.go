package subject

import (
	"context"
	"fmt"

	"backend/internal/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubjectRepository struct {
	connectionURL string
	pool          *pgxpool.Pool
}

func NewSubjectRepository(connectionURL string) *SubjectRepository {
	return &SubjectRepository{connectionURL: connectionURL}
}

func (r *SubjectRepository) Connect(ctx context.Context) error {
	p, err := pgxpool.New(ctx, r.connectionURL)
	if err != nil {
		return fmt.Errorf("pgxpool new: %w", err)
	}

	r.pool = p

	return nil
}

func (r *SubjectRepository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *SubjectRepository) GetAll(ctx context.Context) ([]entities.Subject, error) {
	query := `SELECT id, slug, name_ru, name_kz FROM subjects ORDER BY slug ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get subjects: %w", err)
	}
	defer rows.Close()

	var subjects []entities.Subject
	for rows.Next() {
		var d dto
		if err := rows.Scan(&d.ID, &d.Slug, &d.NameRu, &d.NameKz); err != nil {
			return nil, err
		}
		subjects = append(subjects, d.toEntity())
	}

	return subjects, nil
}

func (r *SubjectRepository) GetByID(ctx context.Context, id string) (*entities.Subject, error) {
	query := `SELECT id, slug, name_ru, name_kz FROM subjects WHERE id = $1`

	var d dto
	err := r.pool.QueryRow(ctx, query, id).Scan(&d.ID, &d.Slug, &d.NameRu, &d.NameKz)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get subject: %w", err)
	}

	subject := d.toEntity()
	return &subject, nil
}

func (r *SubjectRepository) AddInterest(ctx context.Context, userID, subjectID string) error {
	query := `
		INSERT INTO student_interests (user_id, subject_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, subject_id) DO NOTHING
	`

	_, err := r.pool.Exec(ctx, query, userID, subjectID)
	if err != nil {
		return fmt.Errorf("failed to add interest: %w", err)
	}

	return nil
}

func (r *SubjectRepository) RemoveInterest(ctx context.Context, userID, subjectID string) error {
	query := `DELETE FROM student_interests WHERE user_id = $1 AND subject_id = $2`

	_, err := r.pool.Exec(ctx, query, userID, subjectID)
	if err != nil {
		return fmt.Errorf("failed to remove interest: %w", err)
	}
	return nil
}

func (r *SubjectRepository) GetByUserID(ctx context.Context, userID string) ([]entities.Subject, error) {
	query := `
		SELECT s.id, s.slug, s.name_ru, s.name_kz
		FROM subjects s
		JOIN student_interests si ON s.id = si.subject_id
		WHERE si.user_id = $1
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user interests: %w", err)
	}
	defer rows.Close()

	var subjects []entities.Subject
	for rows.Next() {
		var d dto
		if err := rows.Scan(&d.ID, &d.Slug, &d.NameRu, &d.NameKz); err != nil {
			return nil, err
		}
		subjects = append(subjects, d.toEntity())
	}

	return subjects, nil
}

func (r *SubjectRepository) SetInterests(ctx context.Context, userID string, subjectIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM student_interests WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("failed to clear interests: %w", err)
	}

	if len(subjectIDs) > 0 {
		query := `INSERT INTO student_interests (user_id, subject_id) VALUES ($1, $2)`
		for _, subjectID := range subjectIDs {
			_, err := tx.Exec(ctx, query, userID, subjectID)
			if err != nil {
				return fmt.Errorf("failed to insert interest %s: %w", subjectID, err)
			}
		}
	}

	return tx.Commit(ctx)
}
