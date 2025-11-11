package middleware

import (
	"net/http"

	"backend/internal/domain/models"

	"github.com/gin-gonic/gin"
)

// TeacherAuthMiddleware проверяет, что пользователь является учителем или админом.
func TeacherAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in token"})
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role format in token"})
			c.Abort()
			return
		}

		// Пускаем только учителей или админов
		if role != string(models.RoleTeacher) && role != string(models.RoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: teacher or admin role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
