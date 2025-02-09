package routers

import (
	"net/http"
	"nganterin-cs/injectors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func CompRouters(r *gin.RouterGroup, db *gorm.DB, validate *validator.Validate) {
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "pong",
		})
	})

	agentController := injectors.InitializeAgentController(db, validate)
	chatController := injectors.InitializeChatController(db, validate)
	customerController := injectors.InitializeCustomerController(db, validate)

	AgentRoutes(r, agentController)
	ChatRoutes(r, chatController)
	CustomerRoutes(r, customerController)
}
