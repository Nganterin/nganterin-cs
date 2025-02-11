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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
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

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ws.clientsMu.Lock()
	ws.clients[senderData.UUID] = &dto.Connection{
		Conn: conn,
		UUID: senderData.UUID,
		Type: senderData.Type,
	}
	ws.clientsMu.Unlock()

	go ws.pingConnection(conn, senderData.UUID)

	go ws.HandleMessages(ctx, conn, senderData)

	return nil
}

func (ws *WebSocketServiceImpl) pingConnection(conn *websocket.Conn, uuid string) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.RemoveConnection(nil, uuid)
		conn.Close()
	}()

	for range ticker.C {
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			fmt.Printf("Ping error for client %s: %v\n", uuid, err)
			return
		}
	}
}

func (ws *WebSocketServiceImpl) HandleMessages(ctx *gin.Context, conn *websocket.Conn, senderData dto.ChatSender) {
	defer func() {
		ws.RemoveConnection(ctx, senderData.UUID)
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Unexpected close error for client %s: %v\n", senderData.UUID, err)
			}
			break
		}

		var data dto.Chats
		if err := json.Unmarshal(msg, &data); err != nil {
			errorResponse := fmt.Sprintf(`{"error": "Invalid JSON format: %s"}`, err.Error())
			ws.writeMessage(conn, []byte(errorResponse))
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
			ws.writeMessage(conn, []byte(errorResponse))
			continue
		}
	}
}

func (ws *WebSocketServiceImpl) writeMessage(conn *websocket.Conn, message []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.WriteMessage(websocket.TextMessage, message)
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
	ws.clientsMu.Lock()
	defer ws.clientsMu.Unlock()

	for _, client := range ws.clients {
		if client.Type == dto.Agent {
			if err := ws.writeMessage(client.Conn, message); err != nil {
				fmt.Printf("Error sending message to agent %s: %v\n", client.UUID, err)
				go ws.RemoveConnection(nil, client.UUID)
				continue
			}
		}
	}
	return nil
}

func (ws *WebSocketServiceImpl) SendMessageToCustomer(ctx *gin.Context, customerUUID string, message []byte) *exceptions.Exception {
	ws.clientsMu.Lock()
	defer ws.clientsMu.Unlock()

	if targetConn, ok := ws.clients[customerUUID]; ok && targetConn.Type == dto.Customer {
		if err := ws.writeMessage(targetConn.Conn, message); err != nil {
			go ws.RemoveConnection(nil, customerUUID)
			return exceptions.NewException(http.StatusInternalServerError, err.Error())
		}
		return nil
	}
	return exceptions.NewException(http.StatusNotFound, "Customer not found or disconnected")
}

func (ws *WebSocketServiceImpl) RemoveConnection(ctx *gin.Context, uuid string) {
	ws.clientsMu.Lock()
	if client, exists := ws.clients[uuid]; exists {
		client.Conn.Close()
		delete(ws.clients, uuid)
	}
	ws.clientsMu.Unlock()
}
