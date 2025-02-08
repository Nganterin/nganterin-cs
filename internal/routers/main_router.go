package routers

import (
	"nganterin-cs/internal/injectors"

	publicInjector "nganterin-cs/injectors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func InternalRouters(r *gin.RouterGroup, db *gorm.DB, validate *validator.Validate) {
	internalController := injectors.InitializeAuthController(validate)

	agentController := publicInjector.InitializeAgentController(db, validate)

	AuthRoutes(r, internalController)
	AgentRoutes(r, agentController)
}
