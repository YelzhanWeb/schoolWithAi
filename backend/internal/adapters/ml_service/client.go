package mlservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RecommendationResponse struct {
	UserID               string   `json:"user_id"`
	RecommendedCourseIDs []string `json:"recommended_course_ids"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(url string) *Client {
	return &Client{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) GetRecommendedCourseIDs(userID string) ([]string, error) {
	url := fmt.Sprintf("%s/recommend/%s", c.baseURL, userID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned status: %d", resp.StatusCode)
	}

	var result RecommendationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode ML response: %w", err)
	}

	return result.RecommendedCourseIDs, nil
}
