package http

import (
	"context"
	"net/http"

	"backend/internal/adapters/http/handlers"
	"backend/internal/adapters/http/middleware"
	"backend/internal/adapters/storage"
	"backend/internal/services/auth"
	"backend/internal/services/course"
	"backend/internal/services/subject"
	"backend/internal/services/testing"
	"backend/pkg/jwt"

	_ "backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router         *gin.Engine
	httpServer     *http.Server
	authService    *auth.AuthService
	courseService  *course.CourseService
	subjectService *subject.SubjectService
	uploadService  *storage.MinioStorage
	testService    *testing.TestService
	jwtManager     *jwt.JWTManager
	mlServiceURL   string
}

func NewServer(
	authService *auth.AuthService,
	courseService *course.CourseService,
	subjectService *subject.SubjectService,
	uploadService *storage.MinioStorage,
	testService *testing.TestService,
	jwtSecret string,
	mlServiceURL string,
) *Server {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	s := &Server{
		router:         router,
		authService:    authService,
		courseService:  courseService,
		subjectService: subjectService,
		uploadService:  uploadService,
		testService:    testService,
		jwtManager:     jwt.NewJWTManager(jwtSecret),
		mlServiceURL:   mlServiceURL,
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
		courseHandler := handlers.NewCourseHandler(s.courseService)
		subjectHandler := handlers.NewSubjectHandler(s.subjectService)
		uploadHandler := handlers.NewUploadHandler(s.uploadService)
		testHandler := handlers.NewTestHandler(s.testService)

		api.GET("/subjects", subjectHandler.GetAllSubjects)
		api.GET("/tags", courseHandler.GetTags)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(s.jwtManager))
		{
			protected.POST("/auth/change-password", authHandler.ChangePassword)

			protected.POST("/upload", uploadHandler.UploadFile)

			protected.POST("/courses", courseHandler.CreateCourse)
			protected.PUT("/courses/:id", courseHandler.UpdateCourse)
			protected.POST("/courses/:id/publish", courseHandler.ChangePublishStatus)
			protected.GET("/courses/:id/structure", courseHandler.GetStructure)
			protected.GET("/courses/:id", courseHandler.GetCourse)

			protected.POST("/modules", courseHandler.CreateModule)
			protected.PUT("/modules/:id", courseHandler.UpdateModule)
			protected.DELETE("/modules/:id", courseHandler.DeleteModule)

			protected.POST("/lessons", courseHandler.CreateLesson)
			protected.GET("/lessons/:id", courseHandler.GetLesson)
			protected.PUT("/lessons/:id", courseHandler.UpdateLesson)

			protected.POST("/tests", testHandler.CreateTest)
			protected.GET("/modules/:id/test", testHandler.GetTest)
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
