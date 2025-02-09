package services

import (
	"nganterin-cs/api/agents/dto"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompServices interface {
	Create(ctx *gin.Context, data dto.Agents) *exceptions.Exception
	SignIn(ctx *gin.Context, data dto.Login) (*string, *exceptions.Exception)
}
