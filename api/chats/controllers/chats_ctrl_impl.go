package controllers

import (
	"fmt"
	"net/http"
	"nganterin-cs/api/chats/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type CompControllersImpl struct {
	services services.CompServices
}

func NewCompController(compServices services.CompServices) CompControllers {
	return &CompControllersImpl{
		services: compServices,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *CompControllersImpl) ChatWebSocket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		fmt.Println("Received:", string(msg))

		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}