package services

import (
	"backend/internal/domain/models"
	"backend/internal/ports/repositories"
	"backend/pkg/ml_client"
	"context"
	"errors"
	"log"
)

type RecommendationService struct {
	recommendationRepo repositories.CourseRecommendationRepository
	mlClient           *ml_client.MLClient
}

func NewRecommendationService(
	recommendationRepo repositories.CourseRecommendationRepository,
	mlClient *ml_client.MLClient,
) *RecommendationService {
	return &RecommendationService{
		recommendationRepo: recommendationRepo,
		mlClient:           mlClient,
	}
}

func (s *RecommendationService) GetRecommendations(ctx context.Context, studentID int64, limit int) ([]*models.RecommendationResponse, error) {
	savedRecs, err := s.recommendationRepo.GetByStudent(ctx, studentID, limit)

	if err == nil && len(savedRecs) > 0 {
		responseRecs := make([]*models.RecommendationResponse, 0, len(savedRecs))

		for _, rec := range savedRecs {
			responseRecs = append(responseRecs, &models.RecommendationResponse{
				CourseID:  rec.CourseID,
				Title:     rec.CourseTitle, // <-- Используем title из JOIN
				Score:     rec.Score,
				Algorithm: rec.AlgorithmType,
				Reason:    rec.Reason,
				// Details тут не будет, т.к. мы их не храним в БД
			})
		}
		log.Printf("Returning %d recommendations from DB for student %d", len(responseRecs), studentID)
		return responseRecs, nil
	}
	log.Printf("No fresh recommendations in DB for student %d, fetching from ML", studentID)

	// Иначе запросим у ML-сервиса
	mlRecs, err := s.mlClient.GetHybridRecommendations(studentID, limit)
	if err != nil {
		log.Printf("Error fetching recommendations from ML for student %d: %v", studentID, err)
		return nil, errors.New("failed to get recommendations from ML service")
	}
	log.Printf("Fetched %d recommendations from ML for student %d", len(mlRecs), studentID)

	// Сохраним в БД
	var toSave []*models.CourseRecommendation
	for _, mlRec := range mlRecs {
		toSave = append(toSave, &models.CourseRecommendation{
			StudentID:     studentID,
			CourseID:      mlRec.CourseID,
			Score:         mlRec.Score,
			Reason:        mlRec.Reason,
			AlgorithmType: mlRec.Algorithm,
		})
	}

	if len(toSave) > 0 {
		if err := s.recommendationRepo.Save(ctx, toSave); err != nil {
			log.Printf("Error saving recommendations to DB for student %d: %v", studentID, err)
		} else {
			log.Printf("Saved %d recommendations to DB for student %d", len(toSave), studentID)
		}
	}
	return mlRecs, nil
}

// RefreshRecommendations принудительно обновить рекомендации
func (s *RecommendationService) RefreshRecommendations(ctx context.Context, studentID int64, limit int) ([]*models.RecommendationResponse, error) {
	// Удалить старые
	if err := s.recommendationRepo.DeleteOld(ctx, studentID, 0); err != nil {
		log.Printf("Error deleting old recommendations for student %d: %v", studentID, err)
	} else {
		log.Printf("Deleted old recommendations for student %d before refresh", studentID)
	}

	// Получить новые
	mlRecs, err := s.mlClient.GetHybridRecommendations(studentID, limit)
	if err != nil {
		log.Printf("Error fetching recommendations from ML during refresh for student %d: %v", studentID, err)
		return nil, errors.New("failed to refresh recommendations")
	}
	log.Printf("Fetched %d fresh recommendations from ML for student %d", len(mlRecs), studentID)

	// Сохранить
	var toSave []*models.CourseRecommendation
	for _, mlRec := range mlRecs {
		toSave = append(toSave, &models.CourseRecommendation{
			StudentID:     studentID,
			CourseID:      mlRec.CourseID,
			Score:         mlRec.Score,
			Reason:        mlRec.Reason,
			AlgorithmType: mlRec.Algorithm,
		})
	}

	if len(toSave) > 0 {
		if err := s.recommendationRepo.Save(ctx, toSave); err != nil {
			log.Printf("Error saving refreshed recommendations to DB for student %d: %v", studentID, err)
		} else {
			log.Printf("Saved %d refreshed recommendations to DB for student %d", len(toSave), studentID)
		}
	}
	return mlRecs, nil
}
