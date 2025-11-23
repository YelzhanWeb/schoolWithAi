package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/config"
	"backend/internal/adapters/http"
	ml_client "backend/internal/adapters/ml-client"
	postgre "backend/internal/adapters/postgres"
	"backend/internal/services"
	"backend/pkg/jwt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.JWTSecret == "" || cfg.JWTSecret == "your-secret-key-change-this" {
		log.Fatal("❌ JWT_SECRET must be set and should not use default value")
	}

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Database connection established")

	// Repositories
	userRepo := postgre.NewUserRepository(db)
	courseRepo := postgre.NewCourseRepository(db)
	recommendationRepo := postgre.NewRecommendationRepository(db)
	progressRepo := postgre.NewProgressRepository(db)
	profileRepo := postgre.NewStudentProfileRepository(db)
	teacherStatsRepo := postgre.NewTeacherStatisticsRepository(db) // НОВОЕ

	// JWT & ML Client
	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)
	mlServiceURL := config.GetEnv("ML_SERVICE_URL", "http://localhost:5000")
	mlClient := ml_client.NewMLClient(mlServiceURL)

	// Services
	authService := services.NewAuthService(userRepo, jwtManager)
	courseService := services.NewCourseService(courseRepo)
	recommendationService := services.NewRecommendationService(recommendationRepo, mlClient)
	progressService := services.NewProgressService(progressRepo)
	profileService := services.NewProfileService(profileRepo)
	teacherService := services.NewTeacherService(courseRepo, teacherStatsRepo) // НОВОЕ

	// HTTP Server
	httpServer := http.NewServer(
		authService,
		courseService,
		recommendationService,
		progressService,
		profileService,
		teacherService, // НОВОЕ
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

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			db.Close()
			log.Fatal("Server forced to shutdown")
		}
	}

	log.Println("🛑 Server stopped gracefully")
	db.Close()
	log.Println("✅ Database connection closed")
}
