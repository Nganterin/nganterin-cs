package websockets

import (
	"nganterin-cs/api/chats/dto"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type WebSocketServices interface {
	HandleConnection(ctx *gin.Context, senderData dto.ChatSender) *exceptions.Exception
	SendMessageToAgents(ctx *gin.Context, message []byte) *exceptions.Exception
	SendMessageToCustomer(ctx *gin.Context, customerUUID string, message []byte) *exceptions.Exception
	RemoveConnection(ctx *gin.Context, uuid string)
}