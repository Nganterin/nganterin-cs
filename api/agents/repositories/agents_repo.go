package repositories

import (
	"nganterin-cs/models"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompRepositories interface {
	Create(ctx *gin.Context, tx *gorm.DB, data models.Agents) *exceptions.Exception
	FindByUsername(ctx *gin.Context, tx *gorm.DB, username string) *exceptions.Exception
	FindByUUID(ctx *gin.Context, tx *gorm.DB, uuid string) *exceptions.Exception
	FindAll(ctx *gin.Context, tx *gorm.DB) *exceptions.Exception
	Update(ctx *gin.Context, tx *gorm.DB, data models.Agents) *exceptions.Exception
	Delete(ctx *gin.Context, tx *gorm.DB, data models.Agents) *exceptions.Exception
}
