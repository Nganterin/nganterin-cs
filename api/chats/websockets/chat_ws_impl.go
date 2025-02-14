package websockets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nganterin-cs/api/chats/dto"
	"nganterin-cs/api/chats/repositories"
	"nganterin-cs/api/chats/services"
	"nganterin-cs/models"
	"nganterin-cs/pkg/exceptions"
	"nganterin-cs/pkg/helpers"
	"nganterin-cs/pkg/mapper"
	"sync"
	"time"

	customerRepo "nganterin-cs/api/customers/repositories"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const (
	writeWait      = 10 * time.Second
	pingWait       = 120 * time.Second
	maxMessageSize = 512
)

type WebSocketServiceImpl struct {
	repo         repositories.CompRepositories
	customerRepo customerRepo.CompRepositories
	DB           *gorm.DB
	validate     *validator.Validate
	services     services.CompServices
	clients      map[string]*dto.Connection
	clientsMu    sync.Mutex
	upgrader     websocket.Upgrader
}

func NewWebSocketServices(services services.CompServices, repo repositories.CompRepositories, customerRepo customerRepo.CompRepositories, DB *gorm.DB, validate *validator.Validate) WebSocketServices {
	ws := &WebSocketServiceImpl{
		repo:         repo,
		customerRepo: customerRepo,
		DB:           DB,
		validate:     validate,
		services:     services,
		clients:      make(map[string]*dto.Connection),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	go ws.monitorClients()
	return ws
}

func (ws *WebSocketServiceImpl) HandleConnection(ctx *gin.Context, senderData dto.ChatSender) *exceptions.Exception {
	ws.clientsMu.Lock()
	if existingConn, exists := ws.clients[senderData.UUID]; exists {
		existingConn.Conn.Close()
		delete(ws.clients, senderData.UUID)
	}
	ws.clientsMu.Unlock()

	conn, err := ws.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		errorResponse := fmt.Sprintf(`{"error": "%s"}`, err.Error())
		ws.writeMessage(conn, []byte(errorResponse))
		return nil
	}

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pingWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pingWait))
		return nil
	})

	ws.clientsMu.Lock()
	ws.clients[senderData.UUID] = &dto.Connection{
		Conn:     conn,
		UUID:     senderData.UUID,
		Type:     senderData.Type,
		LastPing: time.Now(),
	}
	ws.clientsMu.Unlock()

	if senderData.Type == dto.Agent {
		ws.HandleSendBehindChat(ctx, conn, senderData)
	} else {
		ws.HandleSendCustomerBehindChat(ctx, conn, senderData)
		ws.HandleFirstConnection(ctx, conn, senderData)
	}

	go ws.HandleMessages(ctx, conn, senderData)

	return nil
}

func (ws *WebSocketServiceImpl) HandleFirstConnection(ctx *gin.Context, conn *websocket.Conn, senderData dto.ChatSender) *exceptions.Exception {
	chats, err := ws.repo.FindAllByCustomerUUID(ctx, ws.DB, senderData.UUID)
	if err != nil {
		return err
	}

	if len(chats) == 0 {
		messages := []dto.Chats{
			{
				Type:         dto.Message,
				CustomerUUID: senderData.UUID,
				IsCSChat:     true,
				Message:      "Welcome to Temenin!",
			},
			{
				Type:         dto.Message,
				CustomerUUID: senderData.UUID,
				IsCSChat:     true,
				Message:      "We're here to assist you with any questions or needs related to Nganterin. Feel free to ask anything! 😊",
			},
		}

		for i := range messages {
			ws.ProcessMessage(ctx, &messages[i])
		}
	}

	return nil
}

func (ws *WebSocketServiceImpl) HandleSendBehindChat(ctx *gin.Context, conn *websocket.Conn, senderData dto.ChatSender) *exceptions.Exception {
	var data []models.Chats

	if senderData.LastMessageUUID == "" || senderData.LastMessageUUID == "null" {
		output, err := ws.repo.FindAll(ctx, ws.DB)
		if err != nil {
			return err
		}

		data = output
	} else {
		output, err := ws.repo.FindAllByLastUUID(ctx, ws.DB, senderData.LastMessageUUID)
		if err != nil {
			return err
		}

		data = output
	}

	for _, d := range data {
		output := mapper.MapChatModelToOutput(d)

		msg, exc := json.Marshal(output)
		if exc != nil {
			return exceptions.NewException(http.StatusBadRequest, exc.Error())
		}

		ws.writeMessage(conn, msg)
	}

	return nil
}

func (ws *WebSocketServiceImpl) HandleSendCustomerBehindChat(ctx *gin.Context, conn *websocket.Conn, senderData dto.ChatSender) *exceptions.Exception {
	var data []models.Chats

	if senderData.LastMessageUUID == "" || senderData.LastMessageUUID == "null" {
		output, err := ws.repo.FindAllByCustomerUUID(ctx, ws.DB, senderData.UUID)
		if err != nil {
			return err
		}

		data = output
	} else {
		output, err := ws.repo.FindAllByLastUUIDAndCustomerUUID(ctx, ws.DB, senderData.LastMessageUUID, senderData.UUID)
		if err != nil {
			return err
		}

		data = output
	}

	for _, d := range data {
		output := mapper.MapChatModelToOutput(d)

		msg, exc := json.Marshal(output)
		if exc != nil {
			return exceptions.NewException(http.StatusBadRequest, exc.Error())
		}

		ws.writeMessage(conn, msg)
	}

	return nil
}

func (ws *WebSocketServiceImpl) HandleMessages(ctx *gin.Context, conn *websocket.Conn, senderData dto.ChatSender) {
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

		if data.Type == dto.Ping {
			if err := ws.HandlePing(ctx, senderData); err != nil {
				errorResponse := fmt.Sprintf(`{"error": "%s"}`, err.Error())
				ws.writeMessage(conn, []byte(errorResponse))
				continue
			}
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

func (ws *WebSocketServiceImpl) HandlePing(ctx *gin.Context, senderData dto.ChatSender) *exceptions.Exception {
	ws.clientsMu.Lock()
	if client, exists := ws.clients[senderData.UUID]; exists {
		client.LastPing = time.Now()
		client.Conn.SetReadDeadline(time.Now().Add(pingWait))
	}
	ws.clientsMu.Unlock()
	return nil
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

	go ws.SendMessageToCustomer(ctx, data.CustomerUUID, msg)
	go ws.SendMessageToAgents(ctx, msg)

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

func (ws *WebSocketServiceImpl) monitorClients() {
	for {
		time.Sleep(pingWait / 2)
		ws.clientsMu.Lock()
		for uuid, client := range ws.clients {
			if time.Since(client.LastPing) > pingWait {
				client.Conn.Close()
				delete(ws.clients, uuid)
				fmt.Printf("Client %s disconnected due to inactivity\n", uuid)
			}
		}
		ws.clientsMu.Unlock()
	}
}
