package seed

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"userManagement/internal/models"
)

func SeedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("email = ?", "admin@example.com").Count(&count)
	if count == 0 {
		// Хэшируем пароль
		hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if errPassword != nil {
			log.Fatalf("Не удалось хэшировать пароль: %v", errPassword)
		}

		// Ищем роль админа
		var role models.Role
		if errRole := db.Where("name = ?", "admin").First(&role).Error; errRole != nil {
			log.Fatalf("Не удалось получить роль админа: %v", errRole)
		}

		admin := models.User{
			Name:         "admin",
			Email:        "admin@example.com",
			PasswordHash: string(hashedPassword),
			RoleID:       role.ID,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("Не удалось создать админа: %v", err)
		}

		log.Println("Админ создан: admin@example.com / admin123")
	} else {
		log.Println("Админ уже существует")
	}
}
