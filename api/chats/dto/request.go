package dto

type Chats struct {
	CustomerUUID string `json:"customer_uuid" validate:"required,uuid4"`
	AgentUUID    string `json:"agent_uuid" validate:"required,uuid4"`
	Message      string `json:"message" validate:"required"`
	IsCSChat     bool `json:"is_cs_chat"`
}
