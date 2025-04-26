package main

// @title           User Management API
// @version         1.0
// @description     This is a RESTful API for managing users and roles.
// @termsOfService  http://example.com/terms/

// @contact.name   Владимир Шипунов
// @contact.email  vladimir@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <your-token>
import (
	"net/http"
	"strings"
	"userManagement/internal/config"
	"userManagement/internal/middleware"
	"userManagement/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "userManagement/docs"
)

func main() {
	// Подключаем БД
	config.InitDB()

	// Создаём Gin-роутер
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}
		middleware.RateLimiter()(c)
	})
	// Подключаем все маршруты
	routes.RegisterAllRoutes(r)

	// Подключаем Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("swagger"), ginSwagger.DocExpansion("none")))
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	// Запуск сервера
	r.Run(":8080")
}
