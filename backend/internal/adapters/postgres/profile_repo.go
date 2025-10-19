package postgre

import (
	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type studentProfileRepository struct {
	db *sql.DB
}

func NewStudentProfileRepository(db *sql.DB) repositories.StudentProfileRepository {
	return &studentProfileRepository{db: db}
}

func (r *studentProfileRepository) Create(ctx context.Context, profile *models.StudentProfile) error {
	// Конвертируем preferences в JSON
	preferencesJSON, err := json.Marshal(profile.Preferences)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	query := `
		INSERT INTO student_profiles 
			(user_id, grade, age_group, interests, learning_style, preferences, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err = r.db.QueryRowContext(
		ctx,
		query,
		profile.UserID,
		profile.Grade,
		profile.AgeGroup,
		pq.Array(profile.Interests),
		profile.LearningStyle,
		preferencesJSON,
		time.Now(),
		time.Now(),
	).Scan(&profile.ID)

	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	return nil
}

func (r *studentProfileRepository) GetByUserID(ctx context.Context, userID int64) (*models.StudentProfile, error) {
	query := `
		SELECT id, user_id, grade, age_group, interests, learning_style, preferences, created_at, updated_at
		FROM student_profiles
		WHERE user_id = $1
	`

	profile := &models.StudentProfile{}
	var preferencesJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Grade,
		&profile.AgeGroup,
		pq.Array(&profile.Interests),
		&profile.LearningStyle,
		&preferencesJSON,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Распарсить preferences
	if len(preferencesJSON) > 0 {
		if err := json.Unmarshal(preferencesJSON, &profile.Preferences); err != nil {
			return nil, fmt.Errorf("failed to unmarshal preferences: %w", err)
		}
	} else {
		profile.Preferences = make(map[string]interface{})
	}

	return profile, nil
}

func (r *studentProfileRepository) Update(ctx context.Context, profile *models.StudentProfile) error {
	preferencesJSON, err := json.Marshal(profile.Preferences)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	query := `
		UPDATE student_profiles
		SET grade = $1, age_group = $2, interests = $3, learning_style = $4, 
		    preferences = $5, updated_at = $6
		WHERE user_id = $7
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		profile.Grade,
		profile.AgeGroup,
		pq.Array(profile.Interests),
		profile.LearningStyle,
		preferencesJSON,
		time.Now(),
		profile.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("profile not found")
	}

	return nil
}

func (r *studentProfileRepository) Exists(ctx context.Context, userID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM student_profiles WHERE user_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check profile existence: %w", err)
	}

	return exists, nil
}
