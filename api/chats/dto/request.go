package dto

type Chats struct {
	CustomerUUID string `json:"customer_uuid" validate:"required,uuid4"`
	AgentUUID    string `json:"agent_uuid" validate:"required,uuid4"`
	Message      string `json:"message" validate:"required"`
	IsCSChat     bool   `json:"is_cs_chat"`
}

type Type string

const (
	Agent    Type = "agent"
	Customer Type = "customer"
)

type ChatSender struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  Type   `json:"type"`
}
