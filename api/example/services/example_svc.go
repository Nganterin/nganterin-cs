package services

import (
	"nganterin-cs/api/example/dto"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompService interface {
	Create(ctx *gin.Context, data dto.ExampleInput) *exceptions.Exception
}
