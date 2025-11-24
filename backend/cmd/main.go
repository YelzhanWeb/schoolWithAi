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
	"backend/internal/services/auth"

	"backend/internal/adapters/postgres/user"

	"backend/pkg/jwt"

	_ "github.com/lib/pq"
)

// @title School Backend API
// @version 1.0
// @description API documentation for the School Backend service.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
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

	log.Println("All repositories connected")

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)
	mlServiceURL := config.GetEnv("ML_SERVICE_URL", "http://localhost:5000")
	// mlClient := ml_client.NewMLClient(mlServiceURL)

	authService := auth.NewAuthService(userRepo, jwtManager)

	httpServer := http.NewServer(
		authService,
		cfg.JWTSecret,
		mlServiceURL,
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
