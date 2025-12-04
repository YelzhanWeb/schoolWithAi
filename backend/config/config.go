package config

import (
	"database/sql"
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	ServerPort string
	JWTSecret  string

	// MINIO
	MinioEndpoint   string
	MinioUser       string
	MinioPassword   string
	MinioUseSSL     bool
	MinioBucketName string
	MinioPublicURL  string

	// SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
}

func LoadConfig() Config {
	return Config{
		DBHost:     GetEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     GetEnv("DB_USER", "admin"),
		DBPassword: GetEnv("DB_PASSWORD", "admin123"),
		DBName:     GetEnv("DB_NAME", "education_platform"),
		DBSSLMode:  GetEnv("DB_SSLMODE", "disable"),
		ServerPort: GetEnv("SERVER_PORT", "8080"),
		JWTSecret:  GetEnv("JWT_SECRET", "secret"),
		// MINIO
		MinioEndpoint:   GetEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioUser:       GetEnv("MINIO_ROOT_USER", "admin"),
		MinioPassword:   GetEnv("MINIO_ROOT_PASSWORD", "admin123"),
		MinioUseSSL:     GetEnv("MINIO_USE_SSL", "false") == "true",
		MinioBucketName: GetEnv("MINIO_BUCKET_NAME", "school-assets"),
		MinioPublicURL:  GetEnv("MINIO_PUBLIC_URL", "http://localhost:9000"),
		// SMTP
		SMTPHost:     GetEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:     GetEnv("SMTP_USER", ""),
		SMTPPassword: GetEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     GetEnv("SMTP_FROM", "School With AI <no-reply@school.com>"),
	}
}

func ConnectDB(cfg Config) (*sql.DB, error) {
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

func GetEnv(key, defaultValue string) string {
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
