// go:build wireinject
// go:build wireinject
//go:build wireinject
// +build wireinject

package injectors

import (
	exampleControllers "nganterin-cs/api/example/controllers"
	exampleRepositories "nganterin-cs/api/example/repositories"
	exampleServices "nganterin-cs/api/example/services"
	
	agentControllers "nganterin-cs/api/agents/controllers"
	agentRepositories "nganterin-cs/api/agents/repositories"
	agentServices "nganterin-cs/api/agents/services"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var exampleFeatureSet = wire.NewSet(
	exampleRepositories.NewComponentRepository,
	exampleServices.NewComponentServices,
	exampleControllers.NewCompController,
)

var agentFeatureSet = wire.NewSet(
	agentRepositories.NewComponentRepository,
	agentServices.NewComponentServices,
	agentControllers.NewCompController,
)

func InitializeExampleController(db *gorm.DB, validate *validator.Validate) exampleControllers.CompControllers {
	wire.Build(exampleFeatureSet)
	return nil
}

func InitializeAgentController(db *gorm.DB, validate *validator.Validate) agentControllers.CompControllers {
	wire.Build(agentFeatureSet)
	return nil
}
