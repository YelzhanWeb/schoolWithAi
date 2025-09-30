package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Course struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ExternalURL string `json:"external_url"`
	MinGrade    int    `json:"min_grade"`
	MaxGrade    int    `json:"max_grade"`
	Difficulty  int    `json:"difficulty"`
	IsExternal  bool   `json:"is_external"`
}

type Question struct {
	ID            int      `json:"id"`
	Text          string   `json:"text"`
	QType         string   `json:"q_type"`
	Options       []string `json:"options,omitempty"`
	CorrectAnswer string   `json:"-"`
}

type SubmitAnswer struct {
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
}

type SubmitRequest struct {
	UserID  int            `json:"user_id"`
	Answers []SubmitAnswer `json:"answers"`
}

type SubmitResponse struct {
	TotalQuestions int    `json:"total_questions"`
	Correct        int    `json:"correct"`
	Points         int    `json:"points"`
	Status         string `json:"status"`
}

var db *pgxpool.Pool

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:pass@localhost:5432/appdb"
	}
	var err error
	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	r.Get("/courses", getCourses)
	r.Get("/courses/{id}/test", getCourseQuestions)
	r.Post("/courses/{id}/submit-test", submitTest)
	r.Get("/leaderboard", getLeaderboard)
	r.Get("/users/{id}/score", getUserScore)

	log.Println("Сервер слушает на :8080")
	http.ListenAndServe(":8080", r)
}

// ===== Handlers =====

// список курсов
func getCourses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(),
		"SELECT id, title, description, external_url, min_grade, max_grade, difficulty, is_external FROM courses")
	if err != nil {
		http.Error(w, "Ошибка получения курсов", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.ExternalURL,
			&c.MinGrade, &c.MaxGrade, &c.Difficulty, &c.IsExternal); err != nil {
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}
		courses = append(courses, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// вопросы курса
func getCourseQuestions(w http.ResponseWriter, r *http.Request) {
	courseID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	rows, err := db.Query(context.Background(),
		"SELECT id, text, q_type, options, correct_answer FROM questions WHERE course_id=$1", courseID)
	if err != nil {
		http.Error(w, "Ошибка получения вопросов", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var q Question
		var optionsJSON []byte
		if err := rows.Scan(&q.ID, &q.Text, &q.QType, &optionsJSON, &q.CorrectAnswer); err != nil {
			http.Error(w, "Ошибка чтения вопросов", http.StatusInternalServerError)
			return
		}
		if optionsJSON != nil {
			json.Unmarshal(optionsJSON, &q.Options)
		}
		questions = append(questions, q)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// отправка ответов + начисление баллов
func submitTest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	courseID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var req SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if len(req.Answers) == 0 {
		http.Error(w, "Нет ответов", http.StatusBadRequest)
		return
	}

	correctCount := 0
	for _, ans := range req.Answers {
		var correct string
		err := db.QueryRow(ctx, "SELECT correct_answer FROM questions WHERE id=$1", ans.QuestionID).Scan(&correct)
		if err != nil {
			continue
		}
		isCorrect := strings.EqualFold(strings.TrimSpace(ans.Answer), strings.TrimSpace(correct))

		// сохраняем ответ
		_, _ = db.Exec(ctx,
			`INSERT INTO user_answers (user_id, question_id, answer, is_correct) VALUES ($1,$2,$3,$4)`,
			req.UserID, ans.QuestionID, ans.Answer, isCorrect,
		)

		if isCorrect {
			correctCount++
		}
	}

	total := len(req.Answers)
	percent := float64(correctCount) / float64(total)
	points := correctCount * 5
	status := "in_progress"
	if percent >= 0.6 {
		points += 20
		status = "completed"
	}

	// обновляем прогресс
	_, _ = db.Exec(ctx, `
		INSERT INTO user_progress (user_id, course_id, status, progress_percent, score, last_updated)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (user_id, course_id) DO UPDATE
		SET status=EXCLUDED.status,
		    progress_percent=EXCLUDED.progress_percent,
		    score=EXCLUDED.score,
		    last_updated=EXCLUDED.last_updated
	`, req.UserID, courseID, status, int(percent*100), points, time.Now())

	// обновляем общий рейтинг
	_, _ = db.Exec(ctx, `
		INSERT INTO user_scores (user_id, total_points)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET total_points = user_scores.total_points + $2
	`, req.UserID, points)

	resp := SubmitResponse{
		TotalQuestions: total,
		Correct:        correctCount,
		Points:         points,
		Status:         status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Вставить к остальному коду из предыдущего шага

type LeaderboardEntry struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	TotalPoints int    `json:"total_points"`
	Rank        int    `json:"rank"`
}

// ===== Handlers =====

// leaderboard top N
func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := db.Query(ctx, `
		SELECT u.id, u.name, s.total_points
		FROM user_scores s
		JOIN users u ON u.id = s.user_id
		ORDER BY s.total_points DESC
		LIMIT 10
	`)
	if err != nil {
		http.Error(w, "Ошибка получения рейтинга", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []LeaderboardEntry
	rank := 1
	for rows.Next() {
		var e LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.Name, &e.TotalPoints); err != nil {
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}
		e.Rank = rank
		rank++
		results = append(results, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// очки конкретного пользователя + его место
func getUserScore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var totalPoints int
	var name string
	err := db.QueryRow(ctx,
		`SELECT u.name, COALESCE(s.total_points,0) 
		 FROM users u 
		 LEFT JOIN user_scores s ON u.id=s.user_id
		 WHERE u.id=$1`, userID).Scan(&name, &totalPoints)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	// найти место пользователя в рейтинге
	var rank int
	err = db.QueryRow(ctx,
		`SELECT COUNT(*)+1 
		 FROM user_scores 
		 WHERE total_points > $1`, totalPoints).Scan(&rank)
	if err != nil {
		rank = -1
	}

	resp := LeaderboardEntry{
		UserID:      userID,
		Name:        name,
		TotalPoints: totalPoints,
		Rank:        rank,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
