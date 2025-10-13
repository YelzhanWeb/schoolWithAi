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

// StudentProfileRepository defines the interface for student profile data access
type StudentProfileRepository interface {
	// Create creates a new student profile
	Create(ctx context.Context, profile *models.StudentProfile) error

	// GetByUserID retrieves student profile by user ID
	GetByUserID(ctx context.Context, userID int64) (*models.StudentProfile, error)

	// Update updates student profile
	Update(ctx context.Context, profile *models.StudentProfile) error

	// Delete deletes student profile
	Delete(ctx context.Context, userID int64) error

	// GetByGrade retrieves students by grade
	GetByGrade(ctx context.Context, grade int) ([]*models.StudentProfile, error)

	// GetByAgeGroup retrieves students by age group
	GetByAgeGroup(ctx context.Context, ageGroup string) ([]*models.StudentProfile, error)
}

// CourseRepository defines the interface for course data access
// type CourseRepository interface {
// 	Create(ctx context.Context, course *models.Course) error
// 	GetByID(ctx context.Context, id int64) (*models.Course, error)
// 	Update(ctx context.Context, course *models.Course) error
// 	Delete(ctx context.Context, id int64) error
// 	List(ctx context.Context, filters CourseFilters) ([]*models.Course, error)
// 	GetByTeacher(ctx context.Context, teacherID int64) ([]*models.Course, error)
// }

// CourseFilters defines filters for course queries
type CourseFilters struct {
	Subject       string
	DifficultyMin int
	DifficultyMax int
	AgeGroup      string
	IsPublished   *bool
	Limit         int
	Offset        int
}

// RecommendationRepository defines the interface for recommendations
// type RecommendationRepository interface {
// 	Create(ctx context.Context, recommendation *models.Recommendation) error
// 	GetByStudentID(ctx context.Context, studentID int64, limit int) ([]*models.Recommendation, error)
// 	MarkAsViewed(ctx context.Context, id int64) error
// 	DeleteOld(ctx context.Context, olderThan int) error // delete recommendations older than N days
// }
