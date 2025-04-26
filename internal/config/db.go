package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"userManagement/internal/models"
	"userManagement/internal/seed"
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
		log.Fatal("Ошибка загрузки .env файла: ")
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД: ", err)
	}

	// Сохраняем подключение в глобальной переменной
	DB = db

	// Автоматическая миграция таблицы users
	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Group{},
		&models.ActivityLog{},
	)
	if err != nil {
		log.Fatal("Ошибка миграции: ", err)
	}

	seed.SeedRoles(DB)
	seed.SeedAdmin(DB)

	fmt.Println("Подключение к базе данных успешно! Миграция завершена.")
}
