package routers

import (
	"nganterin-cs/api/chats/controllers"
	"nganterin-cs/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.RouterGroup, controllers controllers.CompControllers) {
	wsGroup := r.Group("/ws")
	wsGroup.Use(middleware.AuthMiddleware())
	{
		wsGroup.GET("/chat", controllers.ChatWebSocket)
	}
}
