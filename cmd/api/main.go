package main

import (
	"github.com/gin-gonic/gin"
	"userManagement/internal/config"
	"userManagement/internal/handler"
)

func main() {
	config.InintDB()

	// Настройка логирования
	r := gin.New()
	r.Use(gin.Logger())   // лог всех входящих запросов
	r.Use(gin.Recovery()) // перехват паники и предотвращение краша

	// Простой маршрут
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.GET("/health", handler.HealthCheckHandler)

	r.POST("/validate", handler.ValidationHandler)

	r.Run(":8080")
}
