package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"userManagement/internal/dto"
	"userManagement/internal/utils"
)

func Authorize(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем информацию о текущем пользователе, которая была добавлена в контекст JWTAuthMiddleware
		currentUser, exists := c.Get("currentUser")
		if !exists {
			utils.Log.Warn("Попытка доступа без авторизации")
			c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Пользователь не авторизован"})
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if currentUser.(dto.UserInfo).Role.Name == role {
				utils.Log.Infof("Доступ разрешен для пользователя ID=%d с ролью %s",
					currentUser.(dto.UserInfo).ID,
					currentUser.(dto.UserInfo).Role.Name)
				c.Next()
				return
			}
		}

		utils.Log.Warnf("Доступ запрещен для пользователя ID=%d с ролью %s",
			currentUser.(dto.UserInfo).ID,
			currentUser.(dto.UserInfo).Role.Name)
		c.JSON(http.StatusForbidden, dto.ResponseError{Message: "Недостаточно прав"})
		c.Abort()
	}
}
