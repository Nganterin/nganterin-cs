// go:build wireinject
// go:build wireinject
//go:build wireinject
// +build wireinject

package injectors

import (
	agentControllers "nganterin-cs/api/agents/controllers"
	agentRepositories "nganterin-cs/api/agents/repositories"
	agentServices "nganterin-cs/api/agents/services"
	
	chatControllers "nganterin-cs/api/chats/controllers"
	chatRepositories "nganterin-cs/api/chats/repositories"
	chatServices "nganterin-cs/api/chats/services"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var agentFeatureSet = wire.NewSet(
	agentRepositories.NewComponentRepository,
	agentServices.NewComponentServices,
	agentControllers.NewCompController,
)

var chatFeatureSet = wire.NewSet(
	chatRepositories.NewComponentRepository,
	chatServices.NewComponentServices,
	chatControllers.NewCompController,
)

func InitializeAgentController(db *gorm.DB, validate *validator.Validate) agentControllers.CompControllers {
	wire.Build(agentFeatureSet)
	return nil
}

func InitializeChatController(db *gorm.DB, validate *validator.Validate) chatControllers.CompControllers {
	wire.Build(chatFeatureSet)
	return nil
}
