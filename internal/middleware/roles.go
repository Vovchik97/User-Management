package middleware

import (
	"net/http"
	"userManagement/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, handlers.ResponseError{Message: "Пользователь не найден"})
			c.Abort()
			return
		}

		user := currentUser.(handlers.UserInfo)
		if user.Role != role {
			c.JSON(http.StatusForbidden, handlers.ResponseError{Message: "Недостаточно прав доступа"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, handlers.ResponseError{Message: "Пользователь не найден"})
			c.Abort()
			return
		}

		user := currentUser.(handlers.UserInfo)
		for _, role := range roles {
			if user.Role == role {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, handlers.ResponseError{Message: "Недостаточно прав доступа"})
		c.Abort()
	}
}
