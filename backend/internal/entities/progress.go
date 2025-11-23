package entities

import "time"

// Progress - прогресс ученика по ресурсу
type Progress struct {
	ID          int64
	StudentID   int64
	ResourceID  int64
	Status      string // not_started, in_progress, completed
	Score       int    // 0-100
	TimeSpent   int    // секунды
	Attempts    int
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ProgressStats - статистика ученика
type ProgressStats struct {
	TotalResources  int
	CompletedCount  int
	InProgressCount int
	AverageScore    float64
	TotalTimeSpent  int // секунды
	TotalPoints     int
	Level           int
	CompletionRate  float64 // процент
}

// ProgressUpdate - запрос на обновление прогресса
type ProgressUpdate struct {
	StudentID  int64
	ResourceID int64
	Status     string
	Score      int
	TimeSpent  int
	Rating     int // опционально, 1-5 звезд
}
