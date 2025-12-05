package http

import (
	"context"
	"net/http"

	"backend/internal/adapters/http/handlers"
	"backend/internal/adapters/http/handlers/content"
	"backend/internal/adapters/http/middleware"
	"backend/internal/adapters/storage"
	"backend/internal/services/auth"
	"backend/internal/services/course"
	"backend/internal/services/gamification"
	"backend/internal/services/student"
	"backend/internal/services/subject"
	"backend/internal/services/testing"
	"backend/pkg/jwt"

	_ "backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router              *gin.Engine
	httpServer          *http.Server
	authService         *auth.AuthService
	courseService       *course.CourseService
	subjectService      *subject.SubjectService
	uploadService       *storage.MinioStorage
	testService         *testing.TestService
	studentService      *student.StudentService
	gamificationService *gamification.GamificationService
	jwtManager          *jwt.JWTManager
}

func NewServer(
	authService *auth.AuthService,
	courseService *course.CourseService,
	subjectService *subject.SubjectService,
	uploadService *storage.MinioStorage,
	testService *testing.TestService,
	studentService *student.StudentService,
	gService *gamification.GamificationService,
	jwtSecret string,
) *Server {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	s := &Server{
		router:              router,
		authService:         authService,
		courseService:       courseService,
		subjectService:      subjectService,
		uploadService:       uploadService,
		testService:         testService,
		studentService:      studentService,
		gamificationService: gService,
		jwtManager:          jwt.NewJWTManager(jwtSecret),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	s.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})

	api := s.router.Group("/v1")
	{
		authHandler := handlers.NewAuthHandler(s.authService)
		courseHandler := content.NewCourseHandler(s.courseService)
		subjectHandler := handlers.NewSubjectHandler(s.subjectService)
		uploadHandler := handlers.NewUploadHandler(s.uploadService)
		testHandler := handlers.NewTestHandler(s.testService)
		studentHandler := handlers.NewStudentHandler(s.studentService)
		gameHandler := handlers.NewGamificationHandler(s.gamificationService)
		leaderboarHandler := handlers.NewLeaderboardHandler(s.studentService)

		api.GET("/subjects", subjectHandler.GetAllSubjects)
		api.GET("/tags", courseHandler.GetTags)
		api.GET("/gamification/leagues", gameHandler.GetAllLeagues)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(s.jwtManager))
		{

			protected.POST("/upload", uploadHandler.UploadFile)

			protected.POST("/auth/change-password", authHandler.ChangePassword)

			protected.POST("/courses/:id/favorite", courseHandler.ToggleFavorite)
			protected.GET("/courses/favorites", courseHandler.GetFavorites)

			protected.POST("/courses", courseHandler.CreateCourse)
			protected.PUT("/courses/:id", courseHandler.UpdateCourse)
			protected.POST("/courses/:id/publish", courseHandler.ChangePublishStatus)
			protected.DELETE("/courses/:id", courseHandler.DeleteCourse)
			protected.GET("/courses/:id/structure", courseHandler.GetStructure)
			protected.GET("/courses/:id", courseHandler.GetCourse)
			protected.GET("/catalog", courseHandler.GetCatalog)

			protected.GET("/teacher/courses", courseHandler.GetMyCourses)

			protected.POST("/modules", courseHandler.CreateModule)
			protected.PUT("/modules/:id", courseHandler.UpdateModule)
			protected.DELETE("/modules/:id", courseHandler.DeleteModule)

			protected.POST("/lessons", courseHandler.CreateLesson)
			protected.GET("/lessons/:id", courseHandler.GetLesson)
			protected.PUT("/lessons/:id", courseHandler.UpdateLesson)
			protected.DELETE("/lessons/:id", courseHandler.DeleteLesson)

			protected.POST("/tests", testHandler.CreateTest)
			protected.GET("/modules/:id/test", testHandler.GetTest)
			protected.GET("/modules/{id}/test-with-answers", testHandler.GetTestWithAnswer)
			protected.PUT("/tests/:id", testHandler.UpdateTest)
			protected.DELETE("/tests/:id", testHandler.DeleteTest)

			protected.POST("/student/onboarding", studentHandler.CompleteOnboarding)
			protected.GET("/student/dashboard", studentHandler.GetDashboard)
			protected.GET("/student/courses/:id/progress", studentHandler.GetCourseProgress)
			protected.POST("/student/lessons/:id/complete", studentHandler.CompleteLesson)
			protected.POST("/student/tests/submit", studentHandler.SubmitTest)
			protected.GET("/student/my-activity-courses", studentHandler.GetAllMyActivityCourses)
			protected.GET("/student/me", studentHandler.GetMe)

			protected.GET("/leaderboard/weekly", leaderboarHandler.GetWeeklyLeaderboard)
			protected.GET("/v1/leaderboard/global", leaderboarHandler.GetGlobalLeaderboard)
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
