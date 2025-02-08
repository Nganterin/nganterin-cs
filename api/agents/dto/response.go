package dto

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Body    interface{} `json:"body,omitempty"`
}

type AgentOutput struct {
	UUID           string `json:"uuid"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	IsActive       string `json:"is_active"`
}
