package repositories

import (
	"backend/internal/domain/models"
	"context"
)

// UserRepository defines the interface for user data access
// Это ПОРТ - контракт, который должен реализовать adapter
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *models.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (*models.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *models.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id int64) error

	// List retrieves users with pagination
	List(ctx context.Context, limit, offset int) ([]*models.User, error)

	// GetByRole retrieves users by role
	GetByRole(ctx context.Context, role models.UserRole) ([]*models.User, error)
}
