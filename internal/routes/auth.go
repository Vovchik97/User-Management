package routes

import (
	"github.com/gin-gonic/gin"
	"userManagement/internal/handlers"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
}
