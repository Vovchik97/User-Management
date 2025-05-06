package handlers

import (
	"net/http"
	"userManagement/internal/config"
	"userManagement/internal/dto"
	"userManagement/internal/models"
	"userManagement/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetActivityLogs godoc
// @Summary      Получить логи активности
// @Tags         Activity
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} models.ActivityLog
// @Failure      401 {object} dto.ResponseError
// @Router       /users/activity [get]
func GetActivityLogs(c *gin.Context) {
	var logs []models.ActivityLog

	// Получаем логи из базы данных
	if err := config.DB.Order("timestamp desc").Find(&logs).Error; err != nil {
		utils.Log.Errorf("Ошибка получения логов активности: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Не удалось получить логи"})
		return
	}

	utils.Log.Info("Логи активности успешно получены")
	c.JSON(http.StatusOK, logs)
}
