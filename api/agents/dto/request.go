package dto

type Agents struct {
	Username string `json:"username" validate:"required,min=4"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required,oneof=agent"`
}
