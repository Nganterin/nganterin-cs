package repositories

import (
	"nganterin-cs/models"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompRepositoriesImpl struct {
}

func NewComponentRepository() CompRepositories {
	return &CompRepositoriesImpl{}
}

func (r *CompRepositoriesImpl) Create(ctx *gin.Context, tx *gorm.DB, data models.Customers) *exceptions.Exception {
	result := tx.Create(&data)
	if result.Error != nil {
		return exceptions.ParseGormError(result.Error)
	}

	return nil
}
