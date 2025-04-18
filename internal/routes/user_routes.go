package routes

import (
	"github.com/gin-gonic/gin"
	"userManagement/internal/handlers"
)

func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("/", handlers.GetUsers)
		users.POST("/", handlers.CreateUser)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
		users.PATCH("/:id/role", handlers.UpdateUserRole)
	}
}
