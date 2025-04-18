package routes

import (
	"github.com/gin-gonic/gin"
	"userManagement/internal/handlers"
	"userManagement/internal/middleware"
)

func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("/", handlers.GetUsers)
		users.POST("/", handlers.CreateUser)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", middleware.RequireAdmin(), handlers.DeleteUser)
		users.PATCH("/:id/role", handlers.UpdateUserRole)
	}
}
