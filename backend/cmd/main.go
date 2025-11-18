package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/adapters/http"
	postgre "backend/internal/adapters/postgres"
	"backend/internal/domain/services"
	"backend/pkg/jwt"
	"backend/pkg/ml_client"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	migrateFlag := flag.String("migrate", "", "Run migrations: up or down")
	flag.Parse()

	cfg := loadConfig()

	// Проверка JWT_SECRET
	if cfg.JWTSecret == "" || cfg.JWTSecret == "your-secret-key-change-this" {
		log.Fatal("❌ JWT_SECRET must be set and should not use default value")
	}

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Database connection established")

	if *migrateFlag != "" {
		handleMigrations(db, *migrateFlag)
		return
	}

	// Repositories
	userRepo := postgre.NewUserRepository(db)
	courseRepo := postgre.NewCourseRepository(db)
	recommendationRepo := postgre.NewRecommendationRepository(db)
	progressRepo := postgre.NewProgressRepository(db)
	profileRepo := postgre.NewStudentProfileRepository(db)
	teacherStatsRepo := postgre.NewTeacherStatisticsRepository(db) // НОВОЕ

	// JWT & ML Client
	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)
	mlServiceURL := getEnv("ML_SERVICE_URL", "http://localhost:5000")
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

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	ServerPort string
	JWTSecret  string
}

func loadConfig() Config {
	jwtSecret := getEnv("JWT_SECRET", "")

	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "education_platform"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  jwtSecret,
	}
}

func connectDB(cfg Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func handleMigrations(db *sql.DB, command string) {
	migrationsPath := "./migrations"

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate: %v", err)
	}

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("✅ Migrations completed")
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("✅ Rollback completed")
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	var value int
	if _, err := fmt.Sscanf(valueStr, "%d", &value); err != nil {
		return defaultValue
	}
	return value
}
