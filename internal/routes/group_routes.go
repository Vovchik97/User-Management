package routes

import (
	"github.com/gin-gonic/gin"
	"userManagement/internal/handlers"
	"userManagement/internal/middleware"
)

func RegisterGroupRoutes(r *gin.Engine) {
	groups := r.Group("/groups")

	// Создание, обновление и удаление групп — только для админа
	groups.POST("/", middleware.Authorize("admin"), handlers.CreateGroups)
	groups.GET("/", middleware.Authorize("admin"), handlers.GetGroups)
	groups.PUT("/:id", middleware.Authorize("admin"), handlers.UpdateGroup)
	groups.DELETE("/:id", middleware.Authorize("admin"), handlers.DeleteGroup)

	// Управление участниками группы — для админа и модератора
	groups.POST("/:id/users", middleware.Authorize("admin", "moderator"), handlers.AddUserToGroup)
	groups.DELETE("/:id/users/:user_id", middleware.Authorize("admin", "moderator"), handlers.RemoveUserFromGroup)
}
