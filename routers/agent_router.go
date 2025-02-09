package routers

import (
	"nganterin-cs/api/agents/controllers"

	"github.com/gin-gonic/gin"
)

func AgentRoutes(r *gin.RouterGroup, controllers controllers.CompControllers) {
	agentGroup := r.Group("/agent")
	{
		authGroup := agentGroup.Group("/auth")
		{
			authGroup.POST("/signin", controllers.SignIn)
		}
	}
}
