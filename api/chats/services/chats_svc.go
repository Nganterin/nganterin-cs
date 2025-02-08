package services

import (
	"nganterin-cs/api/chats/dto"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompServices interface {
	Create(ctx *gin.Context, data dto.Chats) *exceptions.Exception
}
