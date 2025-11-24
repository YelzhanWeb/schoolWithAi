package handlers

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"

	"backend/internal/entities"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Message string `json:"message" example:"something went wrong"`
}
type AuthService interface {
	Register(ctx context.Context, user *entities.User, password string, avatar *multipart.FileHeader) error
	Login(ctx context.Context, email, password string) (string, error)
	ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, email, oldPassword, newPassword string) error
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Email     string                `form:"email"      binding:"required,email"`
	Password  string                `form:"password"   binding:"required,min=8"`
	Role      string                `form:"role"       binding:"required,oneof=student teacher"`
	FirstName string                `form:"first_name" binding:"required"`
	LastName  string                `form:"last_name"  binding:"required"`
	Avatar    *multipart.FileHeader `form:"avatar"`
}

type RegisterResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	FirstName string `json:"full_name"`
	LastName  string `json:"last_name"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user (student or teacher)
// @Tags auth
// @Accept mpfd
// @Produce json
// @Param avatar formData file false "User Avatar Image"
// @Param email formData string true "User Email"
// @Param password formData string true "User Password (min 8 chars)"
// @Param role formData string true "User Role (student or teacher)"
// @Param first_name formData string true "First Name"
// @Param last_name formData string true "Last Name"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "User already exists"
// @Failure 500 {object} ErrorResponse
// @Router /v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Validation failed: " + err.Error()})
		return
	}

	file, _ := c.FormFile("avatar")

	user, err := entities.NewUser(req.Email, req.Password, req.FirstName, req.LastName, "", entities.UserRole(req.Role))
	if err != nil {
		log.Error().Err(err).Msg("domain entity creation failed")
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	err = h.authService.Register(c.Request.Context(), user, req.Password, file)
	if err != nil {
		if errors.Is(err, entities.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, ErrorResponse{Message: err.Error()})
			return
		}
		log.Error().Err(err).Str("email", req.Email).Msg("register user failed")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      string(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})

	log.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Str("role", string(user.Role)).
		Msg("user registered successfully")
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email" example:"student@example.com"`
	Password string `json:"password" binding:"required"       example:"secret123"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid json"})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		}
		log.Error().Err(err).Str("email", req.Email).Msg("login user failed")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})

	log.Info().Str("email", req.Email).Msg("user logged in successfully")
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for the currently authenticated user
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body ChangePasswordRequest true "Password change payload"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse
// @Router /v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid input: " + err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "unauthorized"})
		return
	}

	uidStr, ok := userID.(string)
	if !ok {
		log.Error().Msg("user_id in context is not a string")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error"})
		return
	}

	err := h.authService.ChangePassword(c.Request.Context(), uidStr, req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
			return
		}
		log.Error().Err(err).Str("user_id", uidStr).Msg("failed to change password")
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Info().Str("user_id", uidStr).Msg("password changed successfully")

	c.Status(http.StatusOK)
}

type ResetPasswordRequest struct {
	Email       string `json:"email"        binding:"required,email"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPassword godoc
// @Summary Reset password (public)
// @Description Change password using email and old password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body ResetPasswordRequest true "Reset password payload"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid input: " + err.Error()})
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), req.Email, req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid email or old password"})
			return
		}
		log.Error().Err(err).Str("email", req.Email).Msg("failed to reset password")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "internal server error"})
		return
	}

	log.Info().Str("email", req.Email).Msg("password reset successfully")
	c.Status(http.StatusOK)
}
