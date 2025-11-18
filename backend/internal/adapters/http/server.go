package http

import (
	"context"
	"net/http"

	"backend/internal/adapters/http/handlers"
	"backend/internal/adapters/http/middleware"
	"backend/internal/domain/services"
	"backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router                *gin.Engine
	httpServer            *http.Server
	authService           *services.AuthService
	courseService         *services.CourseService
	recommendationService *services.RecommendationService
	progressService       *services.ProgressService
	profileService        *services.ProfileService
	teacherService        *services.TeacherService
	jwtManager            *jwt.JWTManager
	mlServiceURL          string
}

func NewServer(
	authService *services.AuthService,
	courseService *services.CourseService,
	recommendationService *services.RecommendationService,
	progressService *services.ProgressService,
	profileService *services.ProfileService,
	teacherService *services.TeacherService,
	jwtSecret string,
	mlServiceURL string,
) *Server {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	s := &Server{
		router:                router,
		authService:           authService,
		courseService:         courseService,
		recommendationService: recommendationService,
		progressService:       progressService,
		profileService:        profileService,
		teacherService:        teacherService,
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

	api := s.router.Group("/api")
	{
		// Auth (public)
		authHandler := handlers.NewAuthHandler(s.authService)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(s.jwtManager), authHandler.Me)
		}

		// Courses (public)
		courseHandler := handlers.NewCourseHandler(s.courseService)
		api.GET("/courses", courseHandler.GetAllCourses)
		api.GET("/courses/:id", courseHandler.GetCourse)
		api.GET("/modules/:id/resources", courseHandler.GetModuleResources)

		// // Resources (public)
		// resourceHandler := handlers.NewResourceHandler(s.courseService)
		// api.GET("/resources/:id", resourceHandler.GetResource)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(s.jwtManager))
		{
			// Recommendations
			recHandler := handlers.NewRecommendationHandler(s.recommendationService)
			protected.GET("/recommendations", recHandler.GetRecommendations)
			protected.POST("/recommendations/refresh", recHandler.RefreshRecommendations)

			// Progress
			progressHandler := handlers.NewProgressHandler(s.progressService)
			protected.POST("/progress", progressHandler.UpdateProgress)
			protected.GET("/progress", progressHandler.GetMyProgress)
			protected.GET("/progress/statistics", progressHandler.GetMyStatistics)

			// Profile
			profileHandler := handlers.NewProfileHandler(s.profileService)
			protected.POST("/profile", profileHandler.CreateProfile)
			protected.GET("/profile", profileHandler.GetMyProfile)
			protected.PUT("/profile", profileHandler.UpdateMyProfile)

			// ========================================
			// TEACHER ROUTES
			// ========================================
			teacherApi := protected.Group("/teacher")
			teacherApi.Use(middleware.TeacherAuthMiddleware())
			{
				teacherHandler := handlers.NewTeacherHandler(s.courseService, s.teacherService)

				// Dashboard & Statistics
				teacherApi.GET("/dashboard", teacherHandler.GetDashboard)
				teacherApi.GET("/courses/:id/statistics", teacherHandler.GetCourseStatistics)
				teacherApi.GET("/courses/:id/students", teacherHandler.GetCourseStudents)

				// Course Management
				teacherApi.GET("/courses", teacherHandler.GetMyCourses)
				teacherApi.POST("/courses", teacherHandler.CreateCourse)
				teacherApi.PUT("/courses/:id", teacherHandler.UpdateCourse)
				teacherApi.DELETE("/courses/:id", teacherHandler.DeleteCourse)
				teacherApi.POST("/courses/:id/publish", teacherHandler.PublishCourse)
				teacherApi.POST("/courses/:id/unpublish", teacherHandler.UnpublishCourse)

				// Module Management
				teacherApi.POST("/modules", teacherHandler.CreateModule)
				teacherApi.PUT("/modules/:id", teacherHandler.UpdateModule)
				teacherApi.DELETE("/modules/:id", teacherHandler.DeleteModule)

				// Resource Management
				teacherApi.POST("/resources", teacherHandler.CreateResource)
				teacherApi.PUT("/resources/:id", teacherHandler.UpdateResource)
				teacherApi.DELETE("/resources/:id", teacherHandler.DeleteResource)
			}
		}
	}
}

func (s *Server) Start(port string) error {
	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: s.router,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
