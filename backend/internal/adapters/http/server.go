package http

import (
	"context"
	"net/http"

	"backend/internal/adapters/http/handlers"
	"backend/internal/adapters/http/middleware"
	"backend/internal/services/auth"
	"backend/pkg/jwt"

	_ "backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router       *gin.Engine
	httpServer   *http.Server
	authService  *auth.AuthService
	jwtManager   *jwt.JWTManager
	mlServiceURL string
}

func NewServer(
	authService *auth.AuthService,
	jwtSecret string,
	mlServiceURL string,
) *Server {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	s := &Server{
		router:       router,
		authService:  authService,
		jwtManager:   jwt.NewJWTManager(jwtSecret),
		mlServiceURL: mlServiceURL,
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

	api := s.router.Group("/api")
	{
		authHandler := handlers.NewAuthHandler(s.authService)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(s.jwtManager))
		{
			protected.POST("/auth/change-password", authHandler.ChangePassword)
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
