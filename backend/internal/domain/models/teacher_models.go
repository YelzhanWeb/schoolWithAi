package models

import "time"

// TeacherStatistics - статистика учителя
type TeacherStatistics struct {
	TotalCourses     int                 `json:"total_courses"`
	PublishedCourses int                 `json:"published_courses"`
	DraftCourses     int                 `json:"draft_courses"`
	TotalStudents    int                 `json:"total_students"`
	TotalModules     int                 `json:"total_modules"`
	TotalResources   int                 `json:"total_resources"`
	CourseStats      []*CourseStatistics `json:"course_stats"`
	RecentActivity   []*TeacherActivity  `json:"recent_activity"`
}

// CourseStatistics - статистика по курсу
type CourseStatistics struct {
	CourseID         int64   `json:"course_id"`
	CourseTitle      string  `json:"course_title"`
	EnrolledStudents int     `json:"enrolled_students"`
	CompletedCount   int     `json:"completed_count"`
	InProgressCount  int     `json:"in_progress_count"`
	AverageScore     float64 `json:"average_score"`
	CompletionRate   float64 `json:"completion_rate"`
	AverageRating    float64 `json:"average_rating"`
	TotalRatings     int     `json:"total_ratings"`
}

// TeacherActivity - активность учителя
type TeacherActivity struct {
	ID           int64     `json:"id"`
	TeacherID    int64     `json:"teacher_id"`
	ActivityType string    `json:"activity_type"` // course_created, module_added, resource_added, course_published
	EntityType   string    `json:"entity_type"`   // course, module, resource
	EntityID     int64     `json:"entity_id"`
	EntityTitle  string    `json:"entity_title"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

// StudentProgress - детальный прогресс студента по курсу
type StudentCourseProgress struct {
	StudentID        int64     `json:"student_id"`
	StudentName      string    `json:"student_name"`
	CourseID         int64     `json:"course_id"`
	TotalResources   int       `json:"total_resources"`
	CompletedCount   int       `json:"completed_count"`
	InProgressCount  int       `json:"in_progress_count"`
	AverageScore     float64   `json:"average_score"`
	TotalTimeSpent   int       `json:"total_time_spent"` // в секундах
	LastActivityDate time.Time `json:"last_activity_date"`
	CompletionRate   float64   `json:"completion_rate"`
}

// CreateModuleRequest - запрос на создание модуля
type CreateModuleRequest struct {
	CourseID    int64  `json:"course_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	OrderIndex  int    `json:"order_index" binding:"required"`
}

// CreateResourceRequest - запрос на создание ресурса
type CreateResourceRequest struct {
	ModuleID      int64  `json:"module_id" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Content       string `json:"content"`
	ResourceType  string `json:"resource_type" binding:"required,oneof=exercise quiz reading video interactive"`
	Difficulty    int    `json:"difficulty" binding:"required,min=1,max=5"`
	EstimatedTime int    `json:"estimated_time" binding:"required"` // в минутах
	FileURL       string `json:"file_url"`
	ThumbnailURL  string `json:"thumbnail_url"`
}

// UpdateResourceRequest - запрос на обновление ресурса
type UpdateResourceRequest struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	ResourceType  string `json:"resource_type" binding:"omitempty,oneof=exercise quiz reading video interactive"`
	Difficulty    int    `json:"difficulty" binding:"omitempty,min=1,max=5"`
	EstimatedTime int    `json:"estimated_time"`
	FileURL       string `json:"file_url"`
	ThumbnailURL  string `json:"thumbnail_url"`
}
