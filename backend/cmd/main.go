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
	"backend/internal/adapters/http"
	ml_client "backend/internal/adapters/ml-client"

	"backend/internal/adapters/postgres/course"
	"backend/internal/adapters/postgres/profile"
	"backend/internal/adapters/postgres/subject"
	"backend/internal/adapters/postgres/user"

	"backend/internal/services"
	"backend/pkg/jwt"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.JWTSecret == "" || cfg.JWTSecret == "your-secret-key-change-this" {
		log.Fatal("❌ JWT_SECRET must be set and should not use default value")
	}

	ctx := context.Background()

	// Формируем URL для pgxpool
	connectionURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// ========================================
	// Инициализация репозиториев
	// ========================================

	userRepo := user.NewUserRepository(connectionURL)
	if err := userRepo.Connect(ctx); err != nil {
		log.Fatalf("❌ Failed to connect user repo: %v", err)
	}
	defer userRepo.Close()

	courseRepo := course.NewCourseRepository(connectionURL)
	if err := courseRepo.Connect(ctx); err != nil {
		log.Fatalf("❌ Failed to connect course repo: %v", err)
	}
	defer courseRepo.Close()

	profileRepo := profile.NewStudentProfileRepository(connectionURL)
	if err := profileRepo.Connect(ctx); err != nil {
		log.Fatalf("❌ Failed to connect profile repo: %v", err)
	}
	defer profileRepo.Close()

	subjectRepo := subject.NewSubjectRepository(connectionURL)
	if err := subjectRepo.Connect(ctx); err != nil {
		log.Fatalf("❌ Failed to connect subject repo: %v", err)
	}
	defer subjectRepo.Close()

	log.Println("✅ All repositories connected")

	// ========================================
	// JWT & ML Client
	// ========================================

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)
	mlServiceURL := config.GetEnv("ML_SERVICE_URL", "http://localhost:5000")
	mlClient := ml_client.NewMLClient(mlServiceURL)

	// ========================================
	// Services
	// ========================================

	authService := services.NewAuthService(userRepo, jwtManager)
	courseService := services.NewCourseService(courseRepo)
	profileService := services.NewProfileService(profileRepo, subjectRepo)
	// TODO: recommendationService, progressService, teacherService

	// ========================================
	// HTTP Server
	// ========================================

	httpServer := http.NewServer(
		authService,
		courseService,
		nil, // recommendationService
		nil, // progressService
		profileService,
		nil, // teacherService
		cfg.JWTSecret,
		mlServiceURL,
	)

	log.Println("🚀 Starting Education Platform API...")
	log.Printf("📊 Server running on http://localhost:%s", cfg.ServerPort)

	// Graceful shutdown
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

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			log.Fatal("Server forced to shutdown")
		}
	}

	log.Println("🛑 Server stopped gracefully")
	log.Println("✅ All connections closed")
}
