package http

import (
	"backend/internal/adapters/http/handlers"
	"backend/internal/adapters/http/middleware"
	"backend/internal/domain/services"
	"backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router                *gin.Engine
	authService           *services.AuthService
	courseService         *services.CourseService
	recommendationService *services.RecommendationService // ← НОВОЕ
	jwtManager            *jwt.JWTManager
	mlServiceURL          string
}

func NewServer(
	authService *services.AuthService,
	courseService *services.CourseService,
	recommendationService *services.RecommendationService,
	jwtSecret string,
	mlServiceURL string,
) *Server {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	s := &Server{
		router:                router,
		authService:           authService,
		courseService:         courseService, // ← НОВОЕ
		recommendationService: recommendationService,
		jwtManager:            jwt.NewJWTManager(jwtSecret),
		mlServiceURL:          mlServiceURL,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API группа
	api := s.router.Group("/api")
	{
		// Auth endpoints (без middleware)
		authHandler := handlers.NewAuthHandler(s.authService)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(s.jwtManager), authHandler.Me)
		}

		// Recommendations (protected)
		recHandler := handlers.NewRecommendationHandler(s.recommendationService)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(s.jwtManager))
		{
			protected.GET("/recommendations", recHandler.GetRecommendations)
			protected.POST("/recommendations/refresh", recHandler.RefreshRecommendations)
		}

		courseHandler := handlers.NewCourseHandler(s.courseService)
		api.GET("/courses", courseHandler.GetAllCourses)
		api.GET("/courses/:id", courseHandler.GetCourse)
		api.GET("/modules/:id/resources", courseHandler.GetModuleResources)

	}
}

func (s *Server) Start(port string) error {
	return s.router.Run(":" + port)
}
