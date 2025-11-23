package repositories

import (
	"context"

	models "backend/internal/entities"
)

// StudentProfileRepository - интерфейс для работы с профилем студента
type StudentProfileRepository interface {
	// Создать профиль
	Create(ctx context.Context, profile *models.StudentProfile) error

	// Получить профиль по user_id
	GetByUserID(ctx context.Context, userID int64) (*models.StudentProfile, error)

	// Обновить профиль
	Update(ctx context.Context, profile *models.StudentProfile) error

	// Проверить существование профиля
	Exists(ctx context.Context, userID int64) (bool, error)
}
