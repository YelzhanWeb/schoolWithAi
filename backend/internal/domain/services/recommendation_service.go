package services

import (
	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
	"backend/pkg/ml_client"
	"context"
	"errors"
)

type RecommendationService struct {
	recommendationRepo repositories.RecommendationRepository
	mlClient           *ml_client.MLClient
}

func NewRecommendationService(
	recommendationRepo repositories.RecommendationRepository,
	mlClient *ml_client.MLClient,
) *RecommendationService {
	return &RecommendationService{
		recommendationRepo: recommendationRepo,
		mlClient:           mlClient,
	}
}

// GetRecommendations получить рекомендации из БД или ML-сервиса
func (s *RecommendationService) GetRecommendations(ctx context.Context, studentID int64, limit int) ([]*models.RecommendationResponse, error) {
	// Сначала проверим БД
	savedRecs, err := s.recommendationRepo.GetByStudent(ctx, studentID, limit)

	// Если есть свежие рекомендации (< 1 дня), вернём их
	if err == nil && len(savedRecs) > 0 {
		// TODO: Convert to response format
		return nil, nil
	}

	// Иначе запросим у ML-сервиса
	mlRecs, err := s.mlClient.GetHybridRecommendations(studentID, limit)
	if err != nil {
		return nil, errors.New("failed to get recommendations from ML service")
	}

	// Сохраним в БД
	var toSave []*models.Recommendation
	for _, mlRec := range mlRecs {
		toSave = append(toSave, &models.Recommendation{
			StudentID:     studentID,
			ResourceID:    int64(mlRec.ResourceID),
			Score:         mlRec.Score,
			Reason:        mlRec.Reason,
			AlgorithmType: mlRec.Algorithm,
		})
	}

	_ = s.recommendationRepo.Save(ctx, toSave)

	return mlRecs, nil
}

// RefreshRecommendations принудительно обновить рекомендации
func (s *RecommendationService) RefreshRecommendations(ctx context.Context, studentID int64, limit int) ([]*models.RecommendationResponse, error) {
	// Удалить старые
	_ = s.recommendationRepo.DeleteOld(ctx, studentID, 0)

	// Получить новые
	mlRecs, err := s.mlClient.GetHybridRecommendations(studentID, limit)
	if err != nil {
		return nil, errors.New("failed to refresh recommendations")
	}

	// Сохранить
	var toSave []*models.Recommendation
	for _, mlRec := range mlRecs {
		toSave = append(toSave, &models.Recommendation{
			StudentID:     studentID,
			ResourceID:    int64(mlRec.ResourceID),
			Score:         mlRec.Score,
			Reason:        mlRec.Reason,
			AlgorithmType: mlRec.Algorithm,
		})
	}

	_ = s.recommendationRepo.Save(ctx, toSave)

	return mlRecs, nil
}
