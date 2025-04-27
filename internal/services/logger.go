package services

import (
	"time"
	"userManagement/internal/config"
	"userManagement/internal/models"
)

func LogAction(userID uint, action string) error {
	log := models.ActivityLog{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
	}

	if err := config.DB.Create(&log).Error; err != nil {
		return err
	}

	return nil
}
