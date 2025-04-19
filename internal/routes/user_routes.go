package routes

import (
	"github.com/gin-gonic/gin"
	"userManagement/internal/handlers"
	"userManagement/internal/middleware"
)

func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("/", middleware.Authorize("admin", "moderator", "user"), handlers.GetUsers)
		users.POST("/", middleware.Authorize("admin", "moderator", "user"), handlers.CreateUser)
		users.PUT("/:id", middleware.Authorize("admin", "moderator", "user"), handlers.UpdateUser)
		users.DELETE("/:id", middleware.Authorize("admin", "moderator", "user"), handlers.DeleteUser)
		users.PATCH("/:id/role", middleware.Authorize("admin", "moderator", "user"), handlers.UpdateUserRole)
		users.GET("/activity", middleware.Authorize("admin"), handlers.GetActivityLogs)
		users.PATCH("/:id/ban", middleware.Authorize("admin"), handlers.BanUser)
		users.PATCH("/:id/unban", middleware.Authorize("admin"), handlers.UnbanUser)
	}
}
