package controllers

import (
	"fmt"
	"nganterin-cs/api/chats/services"
	"nganterin-cs/api/chats/websockets"
	"nganterin-cs/pkg/helpers"

	"github.com/gin-gonic/gin"
)

type CompControllersImpl struct {
	services  services.CompServices
	websocket websockets.WebSocketServices
}

func NewCompController(compServices services.CompServices, websocket websockets.WebSocketServices) CompControllers {
	return &CompControllersImpl{
		services:  compServices,
		websocket: websocket,
	}
}

func (h *CompControllersImpl) ChatWebSocket(ctx *gin.Context) {
	senderData, err := helpers.GetSenderData(ctx)
	if err != nil {
		fmt.Println("Failed to get sender data:", err)
		return
	}

	if err := h.websocket.HandleConnection(ctx, senderData); err != nil {
		fmt.Println("Failed to handle WebSocket connection:", err)
		return
	}
}
