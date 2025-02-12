package dto

import (
	"time"

	"github.com/gorilla/websocket"
)

type ChatType string

const (
	Ping    ChatType = "ping"
	Message ChatType = "message"
)

type Chats struct {
	CustomerUUID string   `json:"customer_uuid,omitempty" validate:"required,uuid4"`
	AgentUUID    string   `json:"agent_uuid,omitempty" validate:"required,uuid4"`
	Message      string   `json:"message,omitempty" validate:"required"`
	IsCSChat     bool     `json:"is_cs_chat,omitempty"`
	Type         ChatType `json:"type"`
}

type Type string

const (
	Agent    Type = "agent"
	Customer Type = "customer"
)

type ChatSender struct {
	UUID            string `json:"uuid,omitempty"`
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	Type            Type   `json:"type,omitempty"`
	LastMessageUUID string `json:"last_message_uuid"`
}

type Connection struct {
	Conn     *websocket.Conn
	UUID     string
	Type     Type
	LastPing time.Time
}
