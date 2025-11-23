package services

import (
	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
	"context"
	"errors"
)

type ProfileService struct {
	profileRepo repositories.StudentProfileRepository
}

func NewProfileService(profileRepo repositories.StudentProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
	}
}

// CreateProfile - создать профиль студента
func (s *ProfileService) CreateProfile(ctx context.Context, profile *models.StudentProfile) error {
	// Валидация
	if profile.Grade < 1 || profile.Grade > 11 {
		return errors.New("grade must be between 1 and 11")
	}

	if profile.AgeGroup != "junior" && profile.AgeGroup != "middle" && profile.AgeGroup != "senior" {
		return errors.New("invalid age group")
	}

	// Проверить существование
	exists, err := s.profileRepo.Exists(ctx, profile.UserID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("profile already exists")
	}

	return s.profileRepo.Create(ctx, profile)
}

// GetProfile - получить профиль
func (s *ProfileService) GetProfile(ctx context.Context, userID int64) (*models.StudentProfile, error) {
	return s.profileRepo.GetByUserID(ctx, userID)
}

// UpdateProfile - обновить профиль
func (s *ProfileService) UpdateProfile(ctx context.Context, profile *models.StudentProfile) error {
	// Валидация
	if profile.Grade < 1 || profile.Grade > 11 {
		return errors.New("grade must be between 1 and 11")
	}

	if profile.AgeGroup != "junior" && profile.AgeGroup != "middle" && profile.AgeGroup != "senior" {
		return errors.New("invalid age group")
	}

	return s.profileRepo.Update(ctx, profile)
}
