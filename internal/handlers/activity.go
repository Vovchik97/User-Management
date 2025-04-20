package handlers

import (
	"net/http"
	"userManagement/internal/config"
	"userManagement/internal/models"

	"github.com/gin-gonic/gin"
)

// GetActivityLogs godoc
// @Summary      Получить логи активности
// @Tags         Activity
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} models.ActivityLog
// @Failure      401 {object} ResponseError
// @Router       /users/activity [get]
func GetActivityLogs(c *gin.Context) {
	var logs []models.ActivityLog

	// Получаем логи из базы данных
	if err := config.DB.Order("timestamp desc").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Не удалось получить логи"})
		return
	}

	c.JSON(http.StatusOK, logs)
}
