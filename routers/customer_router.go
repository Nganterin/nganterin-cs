package routers

import (
	"nganterin-cs/api/customers/controllers"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(r *gin.RouterGroup, controllers controllers.CompControllers) {
	customerGroup := r.Group("/customer")
	{
		customerGroup.POST("/register", controllers.Create)
	}
}
