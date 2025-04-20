package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"userManagement/internal/handlers"
)

func Authorize(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем информацию о текущем пользователе, которая была добавлена в контекст JWTAuthMiddleware
		currentUser, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, handlers.ResponseError{Message: "Пользователь не авторизован"})
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if currentUser.(handlers.UserInfo).Role == role {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, handlers.ResponseError{Message: "Недостаточно прав"})
		c.Abort()
	}
}
