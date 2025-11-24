package entities

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
	RoleAdmin   UserRole = "admin"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         UserRole
	FirstName    string
	LastName     string
	AvatarURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(email, password, firstName, lastName, avatarUrl string, role UserRole) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password hash is required")
	}
	if utf8.RuneCountInString(password) < 8 {
		return nil, errors.New("password is less than 8")
	}
	if firstName == "" {
		return nil, errors.New("first name is required")
	}
	if lastName == "" {
		return nil, errors.New("last name is required")
	}
	if !isValidRole(role) {
		return nil, errors.New("invalid user role")
	}

	now := time.Now()
	return &User{
		ID:        uuid.NewString(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		AvatarURL: avatarUrl,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func isValidRole(role UserRole) bool {
	switch role {
	case RoleStudent, RoleTeacher, RoleAdmin:
		return true
	default:
		return false
	}
}
