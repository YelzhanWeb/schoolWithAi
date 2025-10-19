package models

import "time"

type Recommendation struct {
	ID            int64
	StudentID     int64
	ResourceID    int64
	Score         float64
	Reason        string
	AlgorithmType string
	IsViewed      bool
	CreatedAt     time.Time
}

type RecommendationResponse struct {
	ResourceID int     `json:"resource_id"`
	Title      string  `json:"title"`
	Score      float64 `json:"score"`
	Algorithm  string  `json:"algorithm"`
	Reason     string  `json:"reason"`
	Details    *struct {
		CollaborativeScore  float64 `json:"collaborative_score"`
		ContentBasedScore   float64 `json:"content_based_score"`
		KnowledgeBasedScore float64 `json:"knowledge_based_score"`
	} `json:"details,omitempty"`
}
