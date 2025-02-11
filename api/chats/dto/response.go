package dto

import "time"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Body    interface{} `json:"body,omitempty"`
}

type ChatOutput struct {
	UUID         string    `json:"uuid"`
	CustomerUUID string    `json:"customer_uuid"`
	Message      string    `json:"message"`
	IsCSChat     bool      `json:"is_cs_chat"`
	CreatedAt    time.Time `json:"created_at"`

	HumanizedCreatedAt string `json:"humanized_created_at"`

	Customer CustomerOutput `json:"customer"`
}

type CustomerOutput struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
