package repositories

import (
	"nganterin-cs/models"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompRepositories interface {
	Create(ctx *gin.Context, tx *gorm.DB, data models.Chats) *exceptions.Exception
	FindByUUID(ctx *gin.Context, tx *gorm.DB, uuid string) (*models.Chats, *exceptions.Exception)
	FindAll(ctx *gin.Context, tx *gorm.DB) ([]models.Chats, *exceptions.Exception)
	FindAllByLastUUID(ctx *gin.Context, tx *gorm.DB, uuid string) ([]models.Chats, *exceptions.Exception)
	FindAllByCustomerUUID(ctx *gin.Context, tx *gorm.DB, uuid string) ([]models.Chats, *exceptions.Exception)
	FindAllByLastUUIDAndCustomerUUID(ctx *gin.Context, tx *gorm.DB, lastUUID string, customerUUID string) ([]models.Chats, *exceptions.Exception)
}
