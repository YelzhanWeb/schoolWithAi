package testing

import (
	"context"
	"fmt"

	"backend/internal/entities"
)

type TestRepository interface {
	CreateTest(ctx context.Context, test *entities.Test) error
	AddQuestion(ctx context.Context, q *entities.Question) error
	AddAnswer(ctx context.Context, a *entities.Answer) error
	GetTestByModuleID(ctx context.Context, moduleID string) (*entities.Test, error)
}

type TestService struct {
	repo TestRepository
}

func NewTestService(repo TestRepository) *TestService {
	return &TestService{repo: repo}
}

// TO DO: Transaction
func (s *TestService) CreateFullTest(ctx context.Context, test *entities.Test) error {
	if err := s.repo.CreateTest(ctx, test); err != nil {
		return fmt.Errorf("create test: %w", err)
	}

	for _, q := range test.Questions {
		q.TestID = test.ID
		if err := s.repo.AddQuestion(ctx, &q); err != nil {
			return fmt.Errorf("add question: %w", err)
		}

		for _, a := range q.Answers {
			a.QuestionID = q.ID
			if err := s.repo.AddAnswer(ctx, &a); err != nil {
				return fmt.Errorf("add answer: %w", err)
			}
		}
	}
	return nil
}

func (s *TestService) GetTestByModule(ctx context.Context, moduleID string) (*entities.Test, error) {
	test, err := s.repo.GetTestByModuleID(ctx, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test by module id: %w", err)
	}
	return test, nil
}
