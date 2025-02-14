package repositories

import (
	"nganterin-cs/models"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompRepositories interface {
	Create(ctx *gin.Context, tx *gorm.DB, data models.Customers) *exceptions.Exception
	FindByUUID(ctx *gin.Context, tx *gorm.DB, uuid string) (*models.Customers, *exceptions.Exception)
}
