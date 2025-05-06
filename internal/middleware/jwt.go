package middleware

import (
	"net/http"
	"strings"
	"userManagement/internal/config"
	"userManagement/internal/dto"
	"userManagement/internal/models"
	"userManagement/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.Log.Warn("Отсутствует токен авторизации")
			c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Отсутствует токен авторизации"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return config.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			utils.Log.Warnf("Невалидный токен: %s", err)
			c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Невалидный токен"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.Log.Warn("Некорректные данные в токене")
			c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Некорректные данные в токене"})
			c.Abort()
			return
		}

		userIDFloat := claims["userID"].(float64)
		if !ok {
			utils.Log.Warn("Некорректный формат токена")
			c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Некорректный формат токена"})
			c.Abort()
			return
		}
		userID := uint(userIDFloat)

		var user models.User
		if err1 := config.DB.Preload("Role").First(&user, userID).Error; err1 != nil {
			utils.Log.Warnf("Пользователь с ID=%d не найден:", userID)
			c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Пользователь не найден"})
			c.Abort()
			return
		}

		c.Set("userID", user.ID)
		c.Set("currentUser", dto.UserInfo{
			ID:   user.ID,
			Role: user.Role,
		})

		utils.Log.Infof("Успешная авторизация пользователя ID=%d с ролью %s", user.ID, user.Role.Name)
		c.Next()
	}
}
