package dto

type Agents struct {
	Username string `json:"username" validate:"required,min=4"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=agent"`
}
