package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nganterin-cs/api/chats/dto"
	"nganterin-cs/api/chats/services"
	"nganterin-cs/pkg/helpers"
	"sync"

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

type Connection struct {
	Conn *websocket.Conn
	UUID string
	Type dto.Type
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients   = make(map[string]*Connection) // Store connections
	clientsMu sync.Mutex                     // Mutex for concurrent access
)

func (h *CompControllersImpl) ChatWebSocket(ctx *gin.Context) {
	conn, exc := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if exc != nil {
		fmt.Println("WebSocket upgrade error:", exc)
		return
	}

	senderData, err := helpers.GetSenderData(ctx)
	if err != nil {
		conn.Close()
		fmt.Println("Failed to get sender data:", err)
		return
	}

	clientsMu.Lock()
	clients[senderData.UUID] = &Connection{Conn: conn, UUID: senderData.UUID, Type: senderData.Type}
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, senderData.UUID)
		clientsMu.Unlock()
		conn.Close()
	}()

	for {
		_, msg, exc := conn.ReadMessage()
		if exc != nil {
			fmt.Println("Read error:", exc)
			break
		}

		var data dto.Chats
		exc = json.Unmarshal(msg, &data)
		if exc != nil {
			errorResponse := fmt.Sprintf(`{"error": "Invalid JSON format: %s"}`, err.Error())
			conn.WriteMessage(websocket.TextMessage, []byte(errorResponse))
			continue
		}

		fmt.Println("Received:", string(msg))

		clientsMu.Lock()
		if senderData.Type == dto.Customer {
			for _, client := range clients {
				if client.Type == dto.Agent {
					client.Conn.WriteMessage(websocket.TextMessage, msg)
				}
			}
		} else if senderData.Type == dto.Agent {
			if targetConn, ok := clients[data.CustomerUUID]; ok && targetConn.Type == dto.Customer {
				targetConn.Conn.WriteMessage(websocket.TextMessage, msg)
			}
		}

		clientsMu.Unlock()
	}
}
