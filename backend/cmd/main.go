package main

import (
	"backend/internal/adapters/http"
	postgre "backend/internal/adapters/postgres"
	"backend/internal/domain/services"
	"backend/pkg/jwt"
	"backend/pkg/ml_client"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	migrateFlag := flag.String("migrate", "", "Run migrations: up or down")
	flag.Parse()

	cfg := loadConfig()

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
	progressRepo := postgre.NewProgressRepository(db)      // НОВОЕ
	profileRepo := postgre.NewStudentProfileRepository(db) // НОВОЕ

	// JWT & ML Client
	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)
	mlServiceURL := getEnv("ML_SERVICE_URL", "http://localhost:5000")
	mlClient := ml_client.NewMLClient(mlServiceURL)

	// Services
	authService := services.NewAuthService(userRepo, jwtManager)
	courseService := services.NewCourseService(courseRepo)
	recommendationService := services.NewRecommendationService(recommendationRepo, mlClient)
	progressService := services.NewProgressService(progressRepo) // НОВОЕ
	profileService := services.NewProfileService(profileRepo)    // НОВОЕ

	// HTTP Server
	httpServer := http.NewServer(
		authService,
		courseService,
		recommendationService,
		progressService, // НОВОЕ
		profileService,  // НОВОЕ
		cfg.JWTSecret,
		mlServiceURL,
	)

	log.Println("🚀 Starting Education Platform API...")
	log.Printf("📊 Server running on http://localhost:%s", cfg.ServerPort)

	go func() {
		if err := httpServer.Start(cfg.ServerPort); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	log.Println("✅ Server stopped gracefully")
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
	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "education_platform"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-this"),
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
