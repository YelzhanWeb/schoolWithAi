package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type StudentProfile struct {
	ID        string
	UserID    string
	Grade     int
	XP        int64
	Level     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewStudentProfile(userID string, grade int) (*StudentProfile, error) {
	if grade < 1 || grade > 11 {
		return nil, errors.New("grade must be between 1 and 11")
	}

	now := time.Now().UTC()

	return &StudentProfile{
		ID:        uuid.NewString(),
		UserID:    userID,
		Grade:     grade,
		XP:        0,
		Level:     1,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (sp *StudentProfile) AddXP(amount int64) {
	if amount < 0 {
		return
	}
	sp.XP += amount
	sp.UpdatedAt = time.Now().UTC()
}

func (sp *StudentProfile) SetGrade(newGrade int) error {
	if newGrade < 1 || newGrade > 11 {
		return errors.New("grade must be between 1 and 11")
	}
	sp.Grade = newGrade
	sp.UpdatedAt = time.Now().UTC()
	return nil
}
