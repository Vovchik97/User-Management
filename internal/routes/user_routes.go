package routes

import (
	"github.com/gin-gonic/gin"
	"userManagement/internal/handlers"
	"userManagement/internal/middleware"
)

func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	users.Use(middleware.JWTAuthMiddleware())
	{
		users.GET("/me", handlers.GetProfile)

		users.GET("/", middleware.Authorize("admin", "moderator"), handlers.GetUsers)
		users.POST("/", middleware.Authorize("admin"), handlers.CreateUser)
		users.PUT("/:id", middleware.Authorize("admin", "moderator"), handlers.UpdateUser)
		users.DELETE("/:id", middleware.Authorize("admin"), handlers.DeleteUser)
		users.PATCH("/:id/role", middleware.Authorize("admin"), handlers.UpdateUserRole)

		users.GET("/activity", middleware.Authorize("admin"), handlers.GetActivityLogs)
		users.PATCH("/:id/ban", middleware.Authorize("admin"), handlers.BanUser)
		users.PATCH("/:id/unban", middleware.Authorize("admin"), handlers.UnbanUser)
	}
}
