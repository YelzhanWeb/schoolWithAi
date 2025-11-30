package testing

import (
	"context"
	"fmt"
	"time"

	"backend/internal/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TestRepository struct {
	connectionURL string
	pool          *pgxpool.Pool
}

func NewTestRepository(connectionURL string) *TestRepository {
	return &TestRepository{connectionURL: connectionURL}
}

func (r *TestRepository) Connect(ctx context.Context) error {
	p, err := pgxpool.New(ctx, r.connectionURL)
	if err != nil {
		return fmt.Errorf("pgxpool new: %w", err)
	}

	r.pool = p

	return nil
}

func (r *TestRepository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *TestRepository) CreateTest(ctx context.Context, test *entities.Test) error {
	query := `INSERT INTO tests (id, module_id, title, passing_score) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, test.ID, test.ModuleID, test.Title, test.PassingScore)
	return err
}

func (r *TestRepository) AddQuestion(ctx context.Context, q *entities.Question) error {
	query := `INSERT INTO questions (id, test_id, text, question_type) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, q.ID, q.TestID, q.Text, q.QuestionType)
	return err
}

func (r *TestRepository) AddAnswer(ctx context.Context, a *entities.Answer) error {
	query := `INSERT INTO answers (id, question_id, text, is_correct) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, a.ID, a.QuestionID, a.Text, a.IsCorrect)
	return err
}

func (r *TestRepository) GetTestByModuleID(ctx context.Context, moduleID string) (*entities.Test, error) {
	var tDTO testDTO
	queryTest := `SELECT id, module_id, title, passing_score FROM tests WHERE module_id = $1`
	err := r.pool.QueryRow(ctx, queryTest, moduleID).Scan(&tDTO.ID, &tDTO.ModuleID, &tDTO.Title, &tDTO.PassingScore)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrNotFound
		}
		return nil, fmt.Errorf("get test: %w", err)
	}

	test := tDTO.toEntity()

	queryQuestions := `SELECT id, test_id, text, question_type FROM questions WHERE test_id = $1`
	rowsQ, err := r.pool.Query(ctx, queryQuestions, test.ID)
	if err != nil {
		return nil, fmt.Errorf("get questions: %w", err)
	}
	defer rowsQ.Close()

	questionsMap := make(map[string]*entities.Question)
	var questionIDs []string

	for rowsQ.Next() {
		var qDTO questionDTO
		if err := rowsQ.Scan(&qDTO.ID, &qDTO.TestID, &qDTO.Text, &qDTO.QuestionType); err != nil {
			return nil, err
		}
		q := qDTO.toEntity()
		questionsMap[q.ID] = &q
		questionIDs = append(questionIDs, q.ID)
	}

	if len(questionIDs) == 0 {
		return test, nil
	}

	queryAnswers := `SELECT id, question_id, text, is_correct FROM answers WHERE question_id = ANY($1)`
	rowsA, err := r.pool.Query(ctx, queryAnswers, questionIDs)
	if err != nil {
		return nil, fmt.Errorf("get answers: %w", err)
	}
	defer rowsA.Close()

	for rowsA.Next() {
		var aDTO answerDTO
		if err := rowsA.Scan(&aDTO.ID, &aDTO.QuestionID, &aDTO.Text, &aDTO.IsCorrect); err != nil {
			return nil, err
		}
		a := aDTO.toEntity()

		if q, ok := questionsMap[a.QuestionID]; ok {
			q.Answers = append(q.Answers, a)
		}
	}

	for _, id := range questionIDs {
		if q, ok := questionsMap[id]; ok {
			test.Questions = append(test.Questions, *q)
		}
	}

	return test, nil
}

func (r *TestRepository) SaveResult(ctx context.Context, res *entities.TestResult) error {
	if res.AttemptDate.IsZero() {
		res.AttemptDate = time.Now()
	}

	query := `
		INSERT INTO test_results (id, user_id, test_id, score, is_passed, attempt_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query, res.ID, res.UserID, res.TestID, res.Score, res.IsPassed, res.AttemptDate)
	return err
}

func (r *TestRepository) GetUserResults(ctx context.Context, userID string) ([]entities.TestResult, error) {
	query := `
		SELECT id, user_id, test_id, score, is_passed, attempt_date 
		FROM test_results 
		WHERE user_id = $1 
		ORDER BY attempt_date DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []entities.TestResult
	for rows.Next() {
		var rDTO resultDTO
		err := rows.Scan(&rDTO.ID, &rDTO.UserID, &rDTO.TestID, &rDTO.Score, &rDTO.IsPassed, &rDTO.AttemptDate)
		if err != nil {
			return nil, err
		}
		results = append(results, entities.TestResult{
			ID:          rDTO.ID,
			UserID:      rDTO.UserID,
			TestID:      rDTO.TestID,
			Score:       rDTO.Score,
			IsPassed:    rDTO.IsPassed,
			AttemptDate: rDTO.AttemptDate,
		})
	}
	return results, nil
}

func (r *TestRepository) UpdateTest(ctx context.Context, test *entities.Test) error {
	query := `UPDATE tests SET title = $2, passing_score = $3 WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, test.ID, test.Title, test.PassingScore)
	if err != nil {
		return fmt.Errorf("update test: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}
	return nil
}

func (r *TestRepository) DeleteTest(ctx context.Context, testID string) error {
	query := `DELETE FROM tests WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, testID)
	if err != nil {
		return fmt.Errorf("delete test: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.ErrNotFound
	}
	return nil
}

func (r *TestRepository) DeleteQuestionsByTestID(ctx context.Context, testID string) error {
	query := `DELETE FROM questions WHERE test_id = $1`
	_, err := r.pool.Exec(ctx, query, testID)
	return err
}
