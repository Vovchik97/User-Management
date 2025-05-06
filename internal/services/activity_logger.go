package services

import (
	"time"
	"userManagement/internal/config"
	"userManagement/internal/models"
	"userManagement/internal/utils"
)

func LogAction(userID uint, action string) error {
	log := models.ActivityLog{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
	}

	if err := config.DB.Create(&log).Error; err != nil {
		utils.Log.Errorf("Ошибка при записи действия пользователя (ID: %d, Action: %s): %v", userID, action, err)
		return err
	}

	utils.Log.Infof("Записано действие пользователя (ID: %d, Action: %s)", userID, action)
	return nil
}
