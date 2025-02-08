package mapper

import (
	"nganterin-cs/api/agents/dto"
	"nganterin-cs/models"

	"github.com/go-viper/mapstructure/v2"
)

func MapAgentInputToModel(data dto.Agents) models.Agents {
	var result models.Agents

	mapstructure.Decode(data, &result)
	return result
}

func MapAgentModelToOutput(data models.Agents) dto.AgentOutput {
	var result dto.AgentOutput

	mapstructure.Decode(data, &result)
	return result
}
