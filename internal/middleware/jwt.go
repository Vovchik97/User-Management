package middleware

import (
	"net/http"
	"strings"
	"userManagement/internal/config"
	"userManagement/internal/handlers"
	"userManagement/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, handlers.ResponseError{Message: "Отсутствует токен авторизации"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return config.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, handlers.ResponseError{Message: "Невалидный токен"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, handlers.ResponseError{Message: "Некорректные данные в токене"})
			c.Abort()
			return
		}

		userID := uint(claims["userID"].(float64))
		var user models.User
		if err1 := config.DB.First(&user, userID).Error; err1 != nil {
			c.JSON(http.StatusUnauthorized, handlers.ResponseError{Message: "Пользователь не найден"})
			c.Abort()
			return
		}

		c.Set("userID", user.ID)
		c.Set("currentUser", handlers.UserInfo{
			ID:   user.ID,
			Role: user.Role,
		})

		c.Next()
	}
}
