package seed

import (
	"fmt"
	"userManagement/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"userManagement/internal/models"
)

func SeedAdmin(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.User{}).Where("email = ?", "admin@example.com").Count(&count).Error; err != nil {
		utils.Log.Errorf("Ошибка при проверке существующего админа: %v", err)
		return fmt.Errorf("Не удалось проверить наличие админа: %w", err)
	}

	if count > 0 {
		utils.Log.Info("Админ уже существует, пропуск создания")
		return nil
	}

	// Хэшируем пароль
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if errPassword != nil {
		utils.Log.Errorf("Ошибка при хешировании пароля: %v", errPassword)
		return fmt.Errorf("Не удалось хэшировать пароль: %w", errPassword)
	}

	// Ищем роль админа
	var role models.Role
	if errRole := db.Where("name = ?", "admin").First(&role).Error; errRole != nil {
		utils.Log.Errorf("Ошибка при получении роли админа: %v", errRole)
		return fmt.Errorf("Не удалось получить роль админа: %w", errRole)
	}

	admin := models.User{
		Name:         "admin",
		Email:        "admin@example.com",
		PasswordHash: string(hashedPassword),
		RoleID:       role.ID,
	}

	if err := db.Create(&admin).Error; err != nil {
		utils.Log.Fatalf("Не удалось создать админа: %v", err)
	}

	utils.Log.Println("Админ создан: admin@example.com / admin123")

	return nil
}
