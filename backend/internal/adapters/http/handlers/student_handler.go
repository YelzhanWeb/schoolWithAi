package handlers

import (
	"net/http"
	"time"

	"backend/internal/services/student"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	service *student.StudentService
}

func NewStudentHandler(service *student.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

// CompleteLesson godoc
// @Summary Mark lesson as completed
// @Tags student
// @Security BearerAuth
// @Produce json
// @Param id path string true "Lesson ID"
// @Success 200 {object} map[string]any
// @Router /v1/student/lessons/{id}/complete [post]
func (h *StudentHandler) CompleteLesson(c *gin.Context) {
	userID := c.GetString("user_id")
	lessonID := c.Param("id")

	_, xp, err := h.service.CompleteLesson(c.Request.Context(), userID, lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to complete lesson"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "completed",
		"xp_gained": xp,
	})
}

type OnboardingRequest struct {
	Grade      int      `json:"grade" binding:"required,min=1,max=11"`
	SubjectIDs []string `json:"subject_ids" binding:"required"`
}

// CompleteOnboarding godoc
// @Summary Create student profile and set interests
// @Tags student
// @Security BearerAuth
// @Accept json
// @Param input body OnboardingRequest true "Onboarding data"
// @Success 200
// @Router /v1/student/onboarding [post]
func (h *StudentHandler) CompleteOnboarding(c *gin.Context) {
	userID := c.GetString("user_id")
	var req OnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	err := h.service.CompleteOnboarding(c.Request.Context(), userID, req.Grade, req.SubjectIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to complete onboarding: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetCourseProgress godoc
// @Summary Get completed lesson IDs for a course
// @Tags student
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} map[string][]string
// @Router /v1/student/courses/{id}/progress [get]
func (h *StudentHandler) GetCourseProgress(c *gin.Context) {
	userID := c.GetString("user_id")
	courseID := c.Param("id")

	ids, err := h.service.GetCourseProgress(c.Request.Context(), userID, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to get progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"completed_lessons": ids})
}

type DashboardResponse struct {
	Profile       *StudentProfileDTO `json:"profile"`
	Interests     []SubjectDTO       `json:"interests"`
	ActiveCourses []ActiveCourseDTO  `json:"active_courses"`
}

type StudentProfileDTO struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	Grade            int        `json:"grade"`
	XP               int64      `json:"xp"`
	Level            int        `json:"level"`
	CurrentLeagueID  int        `json:"current_league_id"`
	WeeklyXP         int64      `json:"weekly_xp"`
	CurrentStreak    int        `json:"current_streak"`
	MaxStreak        int        `json:"max_streak"`
	LastActivityDate *time.Time `json:"last_activity_date"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type SubjectDTO struct {
	ID     string `json:"id"`
	Slug   string `json:"slug"`
	NameRu string `json:"name_ru"`
	NameKz string `json:"name_kz"`
}

type ActiveCourseDTO struct {
	CourseID           string `json:"course_id"`
	Title              string `json:"title"`
	CoverURL           string `json:"cover_url"`
	ProgressPercentage int    `json:"progress_percentage"`
	TotalLessons       int    `json:"total_lessons"`
	CompletedLessons   int    `json:"completed_lessons"`
}

// GetDashboard godoc
// @Summary Get student dashboard info
// @Tags student
// @Security BearerAuth
// @Produce json
// @Success 200 {object} student.DashboardData
// @Router /v1/student/dashboard [get]
func (h *StudentHandler) GetDashboard(c *gin.Context) {
	userID := c.GetString("user_id")

	data, err := h.service.GetDashboardData(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to get dashboard: " + err.Error()})
		return
	}

	// Если профиля нет (data == nil), возвращаем 404 или специальный флаг
	if data == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "profile_not_found"})
		return
	}

	c.JSON(http.StatusOK, toDashboardDTO(data))
}

func toDashboardDTO(data *student.DashboardData) *DashboardResponse {
	// profile
	profile := &StudentProfileDTO{
		ID:               data.Profile.ID,
		UserID:           data.Profile.UserID,
		Grade:            data.Profile.Grade,
		XP:               data.Profile.XP,
		Level:            data.Profile.Level,
		CurrentLeagueID:  data.Profile.CurrentLeagueID,
		WeeklyXP:         data.Profile.WeeklyXP,
		CurrentStreak:    data.Profile.CurrentStreak,
		MaxStreak:        data.Profile.MaxStreak,
		LastActivityDate: data.Profile.LastActivityDate,
		CreatedAt:        data.Profile.CreatedAt,
		UpdatedAt:        data.Profile.UpdatedAt,
	}

	// interests
	interests := make([]SubjectDTO, len(data.Interests))
	for i, s := range data.Interests {
		interests[i] = SubjectDTO{
			ID:     s.ID,
			Slug:   s.Slug,
			NameRu: s.NameRu,
			NameKz: s.NameKz,
		}
	}

	// courses
	courses := make([]ActiveCourseDTO, len(data.ActiveCourses))
	for i, c := range data.ActiveCourses {
		courses[i] = ActiveCourseDTO{
			CourseID:           c.CourseID,
			Title:              c.Title,
			CoverURL:           c.CoverURL,
			ProgressPercentage: c.ProgressPercentage,
			TotalLessons:       c.TotalLessons,
			CompletedLessons:   c.CompletedLessons,
		}
	}

	return &DashboardResponse{
		Profile:       profile,
		Interests:     interests,
		ActiveCourses: courses,
	}
}

// GetAllMyActivityCourses godoc
// @Summary Get all active courses for student
// @Tags student
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string][]ActiveCourseDTO
// @Router /v1/student/my-courses [get]
func (h *StudentHandler) GetAllMyActivityCourses(c *gin.Context) {
	userID := c.GetString("user_id")

	courses, err := h.service.GetStudentActiveCourses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to get courses"})
		return
	}

	// Маппим в DTO
	dtos := make([]ActiveCourseDTO, len(courses))
	for i, c := range courses {
		dtos[i] = ActiveCourseDTO(c)
	}
	c.JSON(http.StatusOK, gin.H{"courses": dtos})
}

// GetMe godoc
// @Summary Get basic student info for header
// @Tags student
// @Security BearerAuth
// @Produce json
// @Success 200 {object} student.StudentHeaderInfo
// @Router /v1/student/me [get]
func (h *StudentHandler) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")

	info, err := h.service.GetStudentHeaderInfo(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to get user info"})
		return
	}

	c.JSON(http.StatusOK, info)
}
