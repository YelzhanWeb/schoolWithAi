// backend/internal/adapters/http/server.go
package http

import (
	"backend/internal/domain/services"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router      *gin.Engine
	authService *services.AuthService
	// ... другие сервисы
}

func NewServer(authService *services.AuthService) *Server {
	router := gin.Default()

	s := &Server{
		router:      router,
		authService: authService,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// api := s.router.Group("/api")
	{
		// auth := api.Group("/auth")
		{
			// auth.POST("/register", s.handleRegister)
			// auth.POST("/login", s.handleLogin)
		}

		// ... другие routes
	}
}

func (s *Server) Start(port string) error {
	return s.router.Run(":" + port)
}
