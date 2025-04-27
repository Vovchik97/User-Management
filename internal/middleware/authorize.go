package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"userManagement/internal/dto"
)

func Authorize(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем информацию о текущем пользователе, которая была добавлена в контекст JWTAuthMiddleware
		currentUser, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Пользователь не авторизован"})
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if currentUser.(dto.UserInfo).Role.Name == role {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, dto.ResponseError{Message: "Недостаточно прав"})
		c.Abort()
	}
}
