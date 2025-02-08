package mapper

import (
	"nganterin-cs/api/chats/dto"
	"nganterin-cs/models"

	"github.com/go-viper/mapstructure/v2"
)

func MapChatInputToModel(data dto.Chats) models.Chats {
	var result models.Chats

	mapstructure.Decode(data, &result)
	return result
}

func MapChatModelToOutput(data models.Chats) dto.ChatOutput {
	var result dto.ChatOutput

	mapstructure.Decode(data, &result)
	return result
}
