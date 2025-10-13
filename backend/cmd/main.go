package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

func main() {
	// Command line flags
	migrateFlag := flag.String("migrate", "", "Run migrations: up or down")
	flag.Parse()

	// Load configuration
	cfg := loadConfig()

	// Connect to database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Database connection established")

	// Handle migrations
	if *migrateFlag != "" {
		handleMigrations(db, *migrateFlag)
		return
	}

	// Initialize adapters (repositories)
	// userRepo := postgre.NewUserRepository(db)
	// courseRepo := postgres.NewCourseRepository(db)
	// recommendationRepo := postgres.NewRecommendationRepository(db)

	// Initialize domain services (use cases)
	// authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	// courseService := services.NewCourseService(courseRepo)
	// recommendationService := services.NewRecommendationService(recommendationRepo, mlClient)

	// Initialize HTTP handlers (controllers)
	// httpServer := http.NewServer(authService, courseService, recommendationService)

	log.Println("🚀 Starting Education Platform API...")
	log.Printf("📊 Server will run on port %s", cfg.ServerPort)

	// TODO: Start HTTP server
	// go httpServer.Start(cfg.ServerPort)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// TODO: Shutdown HTTP server gracefully
	// httpServer.Shutdown(ctx)

	log.Println("✅ Server stopped gracefully")
}

// Config holds application configuration
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

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func handleMigrations(db *sql.DB, command string) {
	// TODO: Implement migration logic using golang-migrate
	log.Printf("Migration command: %s", command)
	log.Println("⚠️  Migration functionality not yet implemented")
}

// Helper functions
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
