package subject

import (
	"context"
	"fmt"

	"backend/internal/entities"
)

type SubjectRepository interface {
	GetAll(ctx context.Context) ([]entities.Subject, error)
}

type SubjectService struct {
	repo SubjectRepository
}

func NewSubjectService(repo SubjectRepository) *SubjectService {
	return &SubjectService{repo: repo}
}

func (s *SubjectService) GetAllSubjects(ctx context.Context) ([]entities.Subject, error) {
	subjects, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get subjects: %w", err)
	}
	return subjects, nil
}
