package helpers

import (
	"net/http"
	"nganterin-cs/api/chats/dto"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

func GetSenderData(c *gin.Context) (dto.ChatSender, *exceptions.Exception) {
	var result dto.ChatSender
	data, _ := c.Get("sender")

	result, ok := data.(dto.ChatSender)
	if !ok {
		return result, exceptions.NewException(http.StatusUnauthorized, exceptions.ErrInvalidTokenStructure)
	}

	return result, nil
}
