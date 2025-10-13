package models

import (
	"errors"
	"time"
)

// UserRole represents user role in the system
type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
	RoleAdmin   UserRole = "admin"
)

// User represents a user in the system (Domain Entity)
// Это чистая бизнес-модель без зависимостей от БД
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         UserRole
	FullName     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser creates a new user with validation
func NewUser(email, passwordHash, fullName string, role UserRole) (*User, error) {
	// Business validation
	if email == "" {
		return nil, errors.New("email is required")
	}
	if passwordHash == "" {
		return nil, errors.New("password hash is required")
	}
	if fullName == "" {
		return nil, errors.New("full name is required")
	}
	if !isValidRole(role) {
		return nil, errors.New("invalid user role")
	}

	now := time.Now()
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		Role:         role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// IsStudent checks if user is a student
func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

// IsTeacher checks if user is a teacher
func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

// IsAdmin checks if user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate activates the user account
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Helper function
func isValidRole(role UserRole) bool {
	switch role {
	case RoleStudent, RoleTeacher, RoleAdmin:
		return true
	default:
		return false
	}
}

// StudentProfile represents student-specific information
type StudentProfile struct {
	ID            int64
	UserID        int64
	Grade         int
	AgeGroup      string
	Interests     []string
	LearningStyle string
	Preferences   map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewStudentProfile creates a new student profile
func NewStudentProfile(userID int64, grade int, ageGroup, learningStyle string) (*StudentProfile, error) {
	if grade < 1 || grade > 11 {
		return nil, errors.New("grade must be between 1 and 11")
	}

	validAgeGroups := map[string]bool{"junior": true, "middle": true, "senior": true}
	if !validAgeGroups[ageGroup] {
		return nil, errors.New("invalid age group")
	}

	now := time.Now()
	return &StudentProfile{
		UserID:        userID,
		Grade:         grade,
		AgeGroup:      ageGroup,
		LearningStyle: learningStyle,
		Interests:     []string{},
		Preferences:   make(map[string]interface{}),
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// AddInterest adds an interest to student profile
func (sp *StudentProfile) AddInterest(interest string) {
	// Check if interest already exists
	for _, i := range sp.Interests {
		if i == interest {
			return
		}
	}
	sp.Interests = append(sp.Interests, interest)
	sp.UpdatedAt = time.Now()
}

// RemoveInterest removes an interest from student profile
func (sp *StudentProfile) RemoveInterest(interest string) {
	for i, v := range sp.Interests {
		if v == interest {
			sp.Interests = append(sp.Interests[:i], sp.Interests[i+1:]...)
			sp.UpdatedAt = time.Now()
			return
		}
	}
}
