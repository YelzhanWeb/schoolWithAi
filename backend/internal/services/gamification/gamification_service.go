package gamification

import (
	"context"

	"backend/internal/entities"
)

type Repository interface {
	GetAllLeagues(ctx context.Context) ([]entities.League, error)
}

type GamificationService struct {
	repo Repository
}

func NewGamificationService(repo Repository) *GamificationService {
	return &GamificationService{repo: repo}
}

func (s *GamificationService) GetAllLeagues(ctx context.Context) ([]entities.League, error) {
	return s.repo.GetAllLeagues(ctx)
}
