package routes

import (
	"github.com/gin-gonic/gin"
	"userManagement/internal/handlers"
	"userManagement/internal/middleware"
)

func RegisterGroupRoutes(r *gin.Engine) {
	groups := r.Group("/groups")
	groups.Use(middleware.JWTAuthMiddleware())
	{
		// Создание, обновление и удаление групп
		groups.POST("/", middleware.Authorize("admin", "moderator"), handlers.CreateGroups)
		groups.GET("/", middleware.Authorize("admin", "moderator"), handlers.GetGroups)
		groups.PUT("/:id", middleware.Authorize("admin", "moderator"), handlers.UpdateGroup)
		groups.DELETE("/:id", middleware.Authorize("admin", "moderator"), handlers.DeleteGroup)

		// Управление участниками группы
		groups.POST("/:id/users", middleware.Authorize("admin", "moderator"), handlers.AddUserToGroup)
		groups.DELETE("/:id/users/:user_id", middleware.Authorize("admin", "moderator"), handlers.RemoveUserFromGroup)
	}
}
