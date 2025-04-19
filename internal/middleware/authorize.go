package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"userManagement/internal/config"
	"userManagement/internal/handlers"
	"userManagement/internal/models"
)

func Authorize(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, handlers.ResponseError{Message: "Не указан ID пользователя"})
			c.Abort()
			return
		}

		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, handlers.ResponseError{Message: "Пользователь не найден"})
			c.Abort()
			return
		}

		// Добавляем userID в контекст, чтобы логгирование работало
		c.Set("userID", user.ID)
		c.Set("currentUser", user)

		for _, role := range allowedRoles {
			if user.Role == role {
				c.Set("currentUser", user)
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, handlers.ResponseError{Message: "Недостаточно прав"})
		c.Abort()
	}
}
