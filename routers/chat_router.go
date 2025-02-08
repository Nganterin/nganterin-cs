package routers

import (
	"nganterin-cs/api/chats/controllers"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.RouterGroup, controllers controllers.CompControllers) {
	wsGroup := r.Group("/ws")
	{
		wsGroup.GET("/chat", controllers.ChatWebSocket)
	}
}
