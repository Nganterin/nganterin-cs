// go:build wireinject
// go:build wireinject
//go:build wireinject
// +build wireinject

package injectors

import (
	agentControllers "nganterin-cs/api/agents/controllers"
	agentRepositories "nganterin-cs/api/agents/repositories"
	agentServices "nganterin-cs/api/agents/services"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var agentFeatureSet = wire.NewSet(
	agentRepositories.NewComponentRepository,
	agentServices.NewComponentServices,
	agentControllers.NewCompController,
)

func InitializeAgentController(db *gorm.DB, validate *validator.Validate) agentControllers.CompControllers {
	wire.Build(agentFeatureSet)
	return nil
}
