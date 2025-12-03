package user

import (
	"context"
	"time"
)

func (r *UserRepository) SaveResetToken(ctx context.Context, email, token string, ttl time.Duration) error {
	query := `
		INSERT INTO password_reset_tokens (email, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) 
		DO UPDATE SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at, created_at = NOW()
	`
	expiresAt := time.Now().UTC().Add(ttl)
	_, err := r.pool.Exec(ctx, query, email, token, expiresAt)
	return err
}

func (r *UserRepository) VerifyResetToken(ctx context.Context, email, token string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM password_reset_tokens 
			WHERE email = $1 AND token = $2 AND expires_at > NOW()
		)
	`
	var valid bool
	err := r.pool.QueryRow(ctx, query, email, token).Scan(&valid)
	return valid, err
}

func (r *UserRepository) DeleteResetToken(ctx context.Context, email string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM password_reset_tokens WHERE email = $1", email)
	return err
}
