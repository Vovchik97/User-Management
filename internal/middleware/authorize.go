package middleware

import (
	"net/http"
	"userManagement/internal/config"
	"userManagement/internal/handlers"
	"userManagement/internal/models"

	"github.com/gin-gonic/gin"
)

// RequireAdmin проверяет, является ли пользователь администратором
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, handlers.ResponseError{Message: "Не указан ID пользователя"})
			c.Abort()
			return
		}

		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, handlers.ResponseError{Message: "Пользователь не найден"})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, handlers.ResponseError{Message: "Недостаточно прав"})
			c.Abort()
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}
