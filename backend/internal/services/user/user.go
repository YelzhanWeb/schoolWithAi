package user

import (
	"context"
	"fmt"

	"backend/internal/entities"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
}

type UserService struct {
	userRepo UserRepository
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, user *entities.User) error {
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
