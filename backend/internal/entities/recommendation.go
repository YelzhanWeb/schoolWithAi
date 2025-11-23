package entities

import "time"

type CourseRecommendation struct {
	ID            int64
	StudentID     int64
	CourseID      int64 // <-- Используем CourseID
	Score         float64
	Reason        string
	AlgorithmType string
	IsViewed      bool
	CreatedAt     time.Time
	// Можно добавить поля из связанной таблицы courses, если нужно
	CourseTitle string `db:"-"` // Пример, если получаем JOIN'ом
}
type RecommendationResponse struct {
	CourseID  int64   `json:"course_id"` // <-- ИЗМЕНЕНО
	Title     string  `json:"title"`
	Score     float64 `json:"score"`
	Algorithm string  `json:"algorithm"`
	Reason    string  `json:"reason"`
	Details   *struct {
		CollaborativeScore  float64 `json:"collaborative_score"`
		ContentBasedScore   float64 `json:"content_based_score"`
		KnowledgeBasedScore float64 `json:"knowledge_based_score"`
		// AllReasons []string `json:"all_reasons,omitempty"` // <-- Пример, если нужно
	} `json:"details,omitempty"`
}
