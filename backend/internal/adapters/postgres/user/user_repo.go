package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"backend/internal/entities"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	connectionURL string
	pool          *pgxpool.Pool
}

func NewUserRepository(connectionURL string) *UserRepository {
	return &UserRepository{connectionURL: connectionURL}
}

func (r *UserRepository) Connect(ctx context.Context) error {
	p, err := pgxpool.New(ctx, r.connectionURL)
	if err != nil {
		return fmt.Errorf("pgxpool new: %w", err)
	}

	r.pool = p

	return nil
}

func (r *UserRepository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	if r.pool == nil {
		return fmt.Errorf("not connected to pool")
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
		user.UpdatedAt = user.CreatedAt
	}

	d := newDTO(user)

	query := `
		INSERT INTO users (id, email, password_hash, role, first_name, last_name, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		d.ID,
		d.Email,
		d.PasswordHash,
		d.Role,
		d.FirstName,
		d.LastName,
		d.AvatarURL,
		d.CreatedAt,
		d.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entities.ErrAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("not connected to pool")
	}

	query := `
		SELECT id, email, password_hash, role, first_name, last_name, avatar_url, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user, err := scan(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "no rows in result set" {
			return nil, entities.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("not connected to pool")
	}

	query := `
        SELECT id, email, password_hash, role, first_name, last_name, avatar_url, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	user, err := scan(r.pool.QueryRow(ctx, query, email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "no rows in result set" {
			return nil, entities.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	if r.pool == nil {
		return fmt.Errorf("not connected to pool")
	}

	user.UpdatedAt = time.Now().UTC()
	d := newDTO(user)

	query := `
        UPDATE users 
        SET email = $2, 
            password_hash = $3, 
            role = $4, 
            first_name = $5, 
            last_name = $6, 
            avatar_url = $7, 
            updated_at = $8
        WHERE id = $1
    `

	tag, err := r.pool.Exec(
		ctx,
		query,
		d.ID,
		d.Email,
		d.PasswordHash,
		d.Role,
		d.FirstName,
		d.LastName,
		d.AvatarURL,
		d.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	if r.pool == nil {
		return fmt.Errorf("not connected to pool")
	}

	query := `DELETE FROM users WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}

	return nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("not connected to pool")
	}

	query := `
        SELECT id, email, password_hash, role, first_name, last_name, avatar_url, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*entities.User

	for rows.Next() {
		user, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

func (r *UserRepository) GetByRole(ctx context.Context, role entities.UserRole) ([]*entities.User, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("not connected to pool")
	}

	query := `
        SELECT id, email, password_hash, role, first_name, last_name, avatar_url, created_at, updated_at
        FROM users
        WHERE role = $1
        ORDER BY created_at DESC
    `

	rows, err := r.pool.Query(ctx, query, string(role))
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	defer rows.Close()

	var users []*entities.User

	for rows.Next() {
		user, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		users = append(users, &user)
	}

	return users, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scan(scanner rowScanner) (entities.User, error) {
	var d dto

	err := scanner.Scan(
		&d.ID,
		&d.Email,
		&d.PasswordHash,
		&d.Role,
		&d.FirstName,
		&d.LastName,
		&d.AvatarURL,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		return entities.User{}, err
	}

	d.CreatedAt = d.CreatedAt.UTC()
	d.UpdatedAt = d.UpdatedAt.UTC()

	return d.toEntity(), nil
}
