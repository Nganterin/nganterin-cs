package services

import (
	"nganterin-cs/api/customers/dto"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompServices interface {
	Create(ctx *gin.Context, data dto.Customers) (*string, *exceptions.Exception)
}