package user

import (
	"time"

	"backend/internal/entities"
)

type dto struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
	FirstName    string
	LastName     string
	AvatarURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func newDTO(user *entities.User) dto {
	return dto{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         string(user.Role),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		AvatarURL:    user.AvatarURL,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func (d *dto) toEntity() entities.User {
	return entities.User{
		ID:           d.ID,
		Email:        d.Email,
		PasswordHash: d.PasswordHash,
		Role:         entities.UserRole(d.Role),
		FirstName:    d.FirstName,
		LastName:     d.LastName,
		AvatarURL:    d.AvatarURL,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
