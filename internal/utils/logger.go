package utils

import (
	"fmt"
	"time"
	"userManagement/internal/config"
	"userManagement/internal/models"
)

func LogAction(userID uint, action string) {
	log := models.ActivityLog{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
	}

	if err := config.DB.Create(&log).Error; err != nil {
		fmt.Println("Ошибка при записи лога:", err)
	}
}
