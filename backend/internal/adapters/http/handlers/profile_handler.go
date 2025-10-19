package handlers

import (
	"backend/internal/domain/models"
	"backend/internal/domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

type CreateProfileRequest struct {
	Grade         int                    `json:"grade" binding:"required,min=1,max=11"`
	AgeGroup      string                 `json:"age_group" binding:"required"`
	Interests     []string               `json:"interests"`
	LearningStyle string                 `json:"learning_style"`
	Preferences   map[string]interface{} `json:"preferences"`
}

// CreateProfile - создать профиль студента
func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile := &models.StudentProfile{
		UserID:        userID.(int64),
		Grade:         req.Grade,
		AgeGroup:      req.AgeGroup,
		Interests:     req.Interests,
		LearningStyle: req.LearningStyle,
		Preferences:   req.Preferences,
	}

	if err := h.profileService.CreateProfile(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Profile created successfully",
		"profile": profile,
	})
}

// GetMyProfile - получить свой профиль
func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateMyProfile - обновить свой профиль
func (h *ProfileHandler) UpdateMyProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile := &models.StudentProfile{
		UserID:        userID.(int64),
		Grade:         req.Grade,
		AgeGroup:      req.AgeGroup,
		Interests:     req.Interests,
		LearningStyle: req.LearningStyle,
		Preferences:   req.Preferences,
	}

	if err := h.profileService.UpdateProfile(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": profile,
	})
}
