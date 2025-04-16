package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidationInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// HealthCheckHandler отвечает на GET /health — проверка доступности API
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ValidationHandler демонстрирует валидацию входных данных
func ValidationHandler(c *gin.Context) {
	var input ValidationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Validation passed!",
		"data":    input,
	})
}
