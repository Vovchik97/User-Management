package seed

import (
	"fmt"
	"userManagement/internal/models"
	"userManagement/internal/utils"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{Name: "admin", Description: "Администратор"},
		{Name: "moderator", Description: "Модератор"},
		{Name: "user", Description: "Пользователь"},
	}

	for _, role := range roles {
		var existing models.Role
		if err := db.Where("Name = ?", role.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err1 := db.Create(&role).Error; err1 != nil {
					utils.Log.Errorf("Не удалось создать роль %s: %v", role.Name, err1)
					return fmt.Errorf("Не удалось создать роль %s: %v", role.Name, err1)
				} else {
					utils.Log.Printf("Создана роль %s", role.Name)
				}
			} else {
				utils.Log.Errorf("Ошибка при получении роли %s: %v", role.Name, err)
				return fmt.Errorf("Не удалось получить роль %s: %w", role.Name, err)
			}
		} else {
			utils.Log.Infof("Роль %s уже существует, пропуск", role.Name)
		}
	}
	return nil
}
