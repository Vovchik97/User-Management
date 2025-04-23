package seed

import (
	"log"
	"userManagement/internal/models"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
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
					log.Printf("Не удалось создать роль %s: %v", role.Name, err)
				} else {
					log.Printf("Создана роль %s", role.Name)
				}
			}
		}
	}
}
