package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/config"
	"backend/internal/adapters/email"
	"backend/internal/adapters/http"
	"backend/internal/adapters/storage"
	"backend/internal/services/auth"
	"backend/internal/services/scheduler"

	"backend/internal/adapters/postgres/course"
	"backend/internal/adapters/postgres/gamification"
	"backend/internal/adapters/postgres/profile"
	"backend/internal/adapters/postgres/progress"
	"backend/internal/adapters/postgres/subject"
	"backend/internal/adapters/postgres/testing"
	"backend/internal/adapters/postgres/user"
	courseService "backend/internal/services/course"
	gamificationService "backend/internal/services/gamification"
	"backend/internal/services/student"
	subjectService "backend/internal/services/subject"
	testService "backend/internal/services/testing"

	"backend/pkg/jwt"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// @title School Backend API
// @version 1.0
// @description API documentation for the School Backend service.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using default values")
	}
	cfg := config.LoadConfig()

	if cfg.JWTSecret == "" || cfg.JWTSecret == "your-secret-key-change-this" {
		log.Fatal("JWT_SECRET must be set and should not use default value")
	}

	ctx := context.Background()

	connectionURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	userRepo := user.NewUserRepository(connectionURL)
	if err := userRepo.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect user repo: %v", err)
	}
	defer userRepo.Close()

	subjectRepo := subject.NewSubjectRepository(connectionURL)
	if err := subjectRepo.Connect(ctx); err != nil {
		log.Fatalf("Failed subject repo: %v", err)
	}

	courseRepo := course.NewCourseRepository(connectionURL)
	if err := courseRepo.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect course repo: %v", err)
	}
	defer courseRepo.Close()

	testRepo := testing.NewTestRepository(connectionURL)
	if err := testRepo.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect course repo: %v", err)
	}
	defer testRepo.Close()

	profileRepo := profile.NewStudentProfileRepository(connectionURL)
	if err := profileRepo.Connect(ctx); err != nil {
		log.Fatalf("Failed profile repo: %v", err)
	}
	defer profileRepo.Close()

	progressRepo := progress.NewProgressRepository(connectionURL)
	if err := progressRepo.Connect(ctx); err != nil {
		log.Fatalf("Failed profile repo: %v", err)
	}
	defer progressRepo.Close()

	gamificationRepo := gamification.NewGamificationRepository(connectionURL)
	if err := gamificationRepo.Connect(ctx); err != nil {
		log.Fatalf("Failed gamification repo: %v", err)
	}
	defer gamificationRepo.Close()

	log.Println("All repositories connected")

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)

	emailService := email.NewGomailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)
	minioStorage, err := storage.NewMinioStorage(
		cfg.MinioEndpoint,
		cfg.MinioUser,
		cfg.MinioPassword,
		cfg.MinioBucketName,
		cfg.MinioPublicURL,
		cfg.MinioUseSSL,
	)
	if err != nil {
		log.Fatalf("Failed to init MinIO: %v", err)
	}

	authService := auth.NewAuthService(userRepo, jwtManager, minioStorage, emailService)
	subjService := subjectService.NewSubjectService(subjectRepo)
	cService := courseService.NewCourseService(courseRepo)
	testService := testService.NewTestService(testRepo)
	studentService := student.NewStudentService(
		profileRepo,
		subjectRepo,
		progressRepo,
		courseRepo,
		gamificationRepo,
		testRepo,
		userRepo,
	)
	gService := gamificationService.NewGamificationService(gamificationRepo)

	weeklyResetService := scheduler.NewWeeklyResetService(profileRepo, gamificationRepo)
	schedulerCtx, cancelScheduler := context.WithCancel(context.Background())

	go weeklyResetService.Start(schedulerCtx)
	log.Println("Weekly reset scheduler started")

	httpServer := http.NewServer(
		authService,
		cService,
		subjService,
		minioStorage,
		testService,
		studentService,
		gService,
		cfg.JWTSecret,
	)

	log.Println("Starting Education Platform API...")
	log.Printf("Server running on http://localhost:%s", cfg.ServerPort)

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- httpServer.Start(cfg.ServerPort)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)
	case sig := <-quit:
		log.Printf("Received signal: %s", sig)

		cancelScheduler()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			log.Fatal("Server forced to shutdown")
		}
	}

	log.Println("Server stopped gracefully")
	log.Println("All connections closed")
}
