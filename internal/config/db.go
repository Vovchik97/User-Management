package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"userManagement/internal/models"
	"userManagement/internal/seed"
	"userManagement/internal/utils"
)

// Глобальные переменные
var (
	DB        *gorm.DB
	JWTSecret []byte
)

// Инициализируем подключение к БД
func InitDB() {
	// Загружаем .ENV файл
	if err := godotenv.Load(); err != nil {
		utils.Log.Fatal("Ошибка загрузки .env файла: ")
	}

	// Загружаем JWT из .ENV
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))

	// Формируем строку подключения к БД с параметрами из .ENV
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Открываем подключение к БД с помощью GORM
	db, errDB := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDB != nil {
		utils.Log.Fatalf("Ошибка подключения к БД: %v", errDB)
	}

	// Сохраняем подключение в глобальной переменной
	DB = db

	// Автоматическая миграция таблицы users
	errDB = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Group{},
		&models.ActivityLog{},
	)
	if errDB != nil {
		utils.Log.Fatalf("Ошибка миграции: %v", errDB)
	}

	if err := seed.SeedRoles(DB); err != nil {
		utils.Log.Fatalf("Ошибка при сидировании ролей: %v", err)
	}

	if err := seed.SeedAdmin(DB); err != nil {
		utils.Log.Fatalf("Ошибка при сидировании админа: %v", err)
	}

	utils.Log.Info("Подключение к базе данных успешно! Миграция завершена.")
}
