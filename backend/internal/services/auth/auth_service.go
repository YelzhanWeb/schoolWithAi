package auth

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"backend/internal/adapters/email"
	"backend/internal/adapters/storage"
	"backend/internal/entities"
	"backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error

	SaveResetToken(ctx context.Context, email, token string, ttl time.Duration) error
	VerifyResetToken(ctx context.Context, email, token string) (bool, error)
	DeleteResetToken(ctx context.Context, email string) error
}

type AuthService struct {
	userRepo     UserRepository
	jwtManager   *jwt.JWTManager
	storage      storage.FileStorage
	emailService email.EmailService
}

func NewAuthService(userRepo UserRepository, jwtManager *jwt.JWTManager, storage storage.FileStorage, emailService email.EmailService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		jwtManager:   jwtManager,
		storage:      storage,
		emailService: emailService,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	user *entities.User,
	password string,
	avatarFile *multipart.FileHeader,
) error {
	var err error
	user.PasswordHash, err = s.hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hashed password: %w", err)
	}

	if avatarFile != nil {
		avatarURL, err := s.storage.UploadFile(ctx, avatarFile, "avatars")
		if err != nil {
			return fmt.Errorf("failed to upload avatar: %w", err)
		}
		user.AvatarURL = avatarURL
	} else {
		user.AvatarURL = s.storage.GetDefaultAvatarURL()
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if !s.verifyPassword(user.PasswordHash, password) {
		return "", entities.ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	code := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

	if err := s.userRepo.SaveResetToken(ctx, email, code, 15*time.Minute); err != nil {
		return err
	}

	go func() {
		if err := s.emailService.SendResetCode(email, code); err != nil {
			fmt.Printf("Failed to send email to %s: %v\n", email, err)
		}
	}()

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	valid, err := s.userRepo.VerifyResetToken(ctx, email, code)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("invalid or expired reset code")
	}

	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	newHash, err := s.hashPassword(newPassword)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	user.PasswordHash = newHash
	user.UpdatedAt = time.Now().UTC()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return s.userRepo.DeleteResetToken(ctx, email)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !s.verifyPassword(user.PasswordHash, oldPassword) {
		return entities.ErrInvalidCredentials
	}

	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	newPasswordHash, err := s.hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = newPasswordHash
	user.UpdatedAt = time.Now().UTC()

	return s.userRepo.Update(ctx, user)
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

func (s *AuthService) generateToken(user *entities.User) (string, error) {
	return s.jwtManager.Generate(user.ID, user.Email, string(user.Role))
}
