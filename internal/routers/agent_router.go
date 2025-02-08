package routers

import (
	"nganterin-cs/api/agents/controllers"

	"github.com/gin-gonic/gin"
)

func AgentRoutes(r *gin.RouterGroup, controllers controllers.CompControllers) {
	agentGroup := r.Group("/agent")
	{
		agentGroup.POST("/create", controllers.Create)
	}
}
