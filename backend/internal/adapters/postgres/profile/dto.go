package profile

import (
	"time"

	"backend/internal/entities"
)

type dto struct {
	ID        string
	UserID    string
	Grade     int
	XP        int64
	Level     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func newDTO(sp *entities.StudentProfile) dto {
	return dto{
		ID:        sp.ID,
		UserID:    sp.UserID,
		Grade:     sp.Grade,
		XP:        sp.XP,
		Level:     sp.Level,
		CreatedAt: sp.CreatedAt,
		UpdatedAt: sp.UpdatedAt,
	}
}

func (d *dto) toEntity() *entities.StudentProfile {
	return &entities.StudentProfile{
		ID:        d.ID,
		UserID:    d.UserID,
		Grade:     d.Grade,
		XP:        d.XP,
		Level:     d.Level,
		CreatedAt: d.CreatedAt.UTC(),
		UpdatedAt: d.UpdatedAt.UTC(),
	}
}
