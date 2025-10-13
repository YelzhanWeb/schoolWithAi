package services

import (
	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
// Это USE CASE - чистая бизнес-логика
type AuthService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// RegisterRequest contains registration data
type RegisterRequest struct {
	Email    string
	Password string
	FullName string
	Role     models.UserRole
}

// LoginRequest contains login credentials
type LoginRequest struct {
	Email    string
	Password string
}

// AuthResponse contains authentication result
type AuthResponse struct {
	User  *models.User
	Token string
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	// Validate input
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	passwordHash, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user domain model
	user, err := models.NewUser(req.Email, passwordHash, req.FullName, req.Role)
	if err != nil {
		return nil, err
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if account is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !s.verifyPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// ChangePassword changes user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !s.verifyPassword(user.PasswordHash, oldPassword) {
		return errors.New("invalid old password")
	}

	// Validate new password
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Hash new password
	newPasswordHash, err := s.hashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update user
	user.PasswordHash = newPasswordHash
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

// DeactivateUser deactivates a user account
func (s *AuthService) DeactivateUser(ctx context.Context, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.Deactivate()
	return s.userRepo.Update(ctx, user)
}

// Private methods

func (s *AuthService) validateRegisterRequest(req RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if req.FullName == "" {
		return errors.New("full name is required")
	}
	return nil
}

func (s *AuthService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *AuthService) verifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	// TODO: Implement JWT token generation
	// This is a placeholder
	return "jwt-token-placeholder", nil
}
