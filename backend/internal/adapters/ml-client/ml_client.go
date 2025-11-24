package ml_client

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"time"

// 	"backend/internal/entities"
// )

// type MLClient struct {
// 	baseURL    string
// 	httpClient *http.Client
// }

// type recommendationRequest struct {
// 	StudentID int64 `json:"student_id"`
// 	TopN      int   `json:"top_n"`
// }

// func NewMLClient(baseURL string) *MLClient {
// 	return &MLClient{
// 		baseURL: baseURL,
// 		httpClient: &http.Client{
// 			Timeout: 10 * time.Second,
// 		},
// 	}
// }

// func (c *MLClient) GetHybridRecommendations(studentID int64, topN int) ([]*entities.RecommendationResponse, error) {
// 	reqBody := recommendationRequest{
// 		StudentID: studentID,
// 		TopN:      topN,
// 	}

// 	jsonData, err := json.Marshal(reqBody)
// 	if err != nil {
// 		return nil, err
// 	}

// 	url := fmt.Sprintf("%s/recommendations/hybrid", c.baseURL)
// 	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to call ML service: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := io.ReadAll(resp.Body)
// 		return nil, fmt.Errorf("ML service error (%d): %s", resp.StatusCode, string(body))
// 	}

// 	var recommendations []*entities.RecommendationResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&recommendations); err != nil {
// 		return nil, fmt.Errorf("failed to decode response: %w", err)
// 	}

// 	return recommendations, nil
// }
