package websockets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nganterin-cs/api/chats/dto"
	"nganterin-cs/api/chats/repositories"
	"nganterin-cs/api/chats/services"
	"nganterin-cs/pkg/exceptions"
	"nganterin-cs/pkg/helpers"
	"nganterin-cs/pkg/mapper"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type WebSocketServiceImpl struct {
	repo      repositories.CompRepositories
	DB        *gorm.DB
	validate  *validator.Validate
	services  services.CompServices
	clients   map[string]*dto.Connection
	clientsMu sync.Mutex
	upgrader  websocket.Upgrader
}

func NewWebSocketServices(services services.CompServices, repo repositories.CompRepositories, DB *gorm.DB, validate *validator.Validate) WebSocketServices {
	return &WebSocketServiceImpl{
		repo:     repo,
		DB:       DB,
		validate: validate,
		services: services,
		clients:  make(map[string]*dto.Connection),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (ws *WebSocketServiceImpl) HandleConnection(ctx *gin.Context, senderData dto.ChatSender) *exceptions.Exception {
	conn, err := ws.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return exceptions.NewException(http.StatusInternalServerError, err.Error())
	}

	ws.clientsMu.Lock()
	ws.clients[senderData.UUID] = &dto.Connection{
		Conn: conn,
		UUID: senderData.UUID,
		Type: senderData.Type,
	}
	ws.clientsMu.Unlock()

	go ws.HandleMessages(ctx, conn, senderData)

	return nil
}

func (ws *WebSocketServiceImpl) HandleMessages(ctx *gin.Context, conn *websocket.Conn, senderData dto.ChatSender) {
	defer func() {
		ws.RemoveConnection(ctx, senderData.UUID)
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		var data dto.Chats
		if err := json.Unmarshal(msg, &data); err != nil {
			errorResponse := fmt.Sprintf(`{"error": "Invalid JSON format: %s"}`, err.Error())
			conn.WriteMessage(websocket.TextMessage, []byte(errorResponse))
			continue
		}

		if senderData.Type == dto.Customer {
			data.CustomerUUID = senderData.UUID
			data.IsCSChat = false
		} else if senderData.Type == dto.Agent {
			data.AgentUUID = senderData.UUID
			data.IsCSChat = true
		}

		if err := ws.ProcessMessage(ctx, &data); err != nil {
			errorResponse := fmt.Sprintf(`{"error": "%s"}`, err.Error())
			conn.WriteMessage(websocket.TextMessage, []byte(errorResponse))
			continue
		}
	}
}

func (ws *WebSocketServiceImpl) ProcessMessage(ctx *gin.Context, data *dto.Chats) *exceptions.Exception {
	ws.clientsMu.Lock()
	defer ws.clientsMu.Unlock()

	input := mapper.MapChatInputToModel(*data)
	input.UUID = uuid.NewString()

	tx := ws.DB.Begin()
	defer helpers.CommitOrRollback(tx)

	err := ws.repo.Create(ctx, tx, input)
	if err != nil {
		return err
	}

	result, err := ws.repo.FindByUUID(ctx, tx, input.UUID)
	if err != nil {
		return err
	}

	output := mapper.MapChatModelToOutput(*result)

	msg, exc := json.Marshal(output)
	if exc != nil {
		return exceptions.NewException(http.StatusBadRequest, exc.Error())
	}

	if data.IsCSChat {
		go ws.SendMessageToCustomer(ctx, data.CustomerUUID, msg)
	} else {
		go ws.SendMessageToAgents(ctx, msg)
	}

	return nil
}

func (ws *WebSocketServiceImpl) SendMessageToAgents(ctx *gin.Context, message []byte) *exceptions.Exception {
	for _, client := range ws.clients {
		if client.Type == dto.Agent {
			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return exceptions.NewException(http.StatusInternalServerError, err.Error())
			}
		}
	}
	return nil
}

func (ws *WebSocketServiceImpl) SendMessageToCustomer(ctx *gin.Context, customerUUID string, message []byte) *exceptions.Exception {
	if targetConn, ok := ws.clients[customerUUID]; ok && targetConn.Type == dto.Customer {
		err := targetConn.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return exceptions.NewException(http.StatusInternalServerError, err.Error())
		}

		return nil
	}
	return exceptions.NewException(http.StatusNotFound, "Customer not found or disconnected")
}

func (ws *WebSocketServiceImpl) RemoveConnection(ctx *gin.Context, uuid string) {
	ws.clientsMu.Lock()
	delete(ws.clients, uuid)
	ws.clientsMu.Unlock()
}
