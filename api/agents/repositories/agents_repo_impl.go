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

func (r *CompRepositoriesImpl) Create(ctx *gin.Context, tx *gorm.DB, data models.Agents) *exceptions.Exception {
	result := tx.Create(&data)
	if result.Error != nil {
		return exceptions.ParseGormError(result.Error)
	}

	return nil
}

func (r *CompRepositoriesImpl) FindByUsername(ctx *gin.Context, tx *gorm.DB, username string) *exceptions.Exception {
	var agent models.Agents

	result := tx.Where("username = ?", username).First(&agent)
	if result.Error != nil {
		return exceptions.ParseGormError(result.Error)
	}

	return nil
}

func (r *CompRepositoriesImpl) FindByUUID(ctx *gin.Context, tx *gorm.DB, uuid string) *exceptions.Exception {
	var agent models.Agents

	result := tx.Where("uuid = ?", uuid).First(&agent)
	if result.Error != nil {
		return exceptions.ParseGormError(result.Error)
	}

	return nil
}

func (r *CompRepositoriesImpl) FindAll(ctx *gin.Context, tx *gorm.DB) *exceptions.Exception {
	var agents []models.Agents

	result := tx.Find(&agents)
	if result.Error != nil {
		return exceptions.ParseGormError(result.Error)
	}

	return nil
}

func (r *CompRepositoriesImpl) Update(ctx *gin.Context, tx *gorm.DB, data models.Agents) *exceptions.Exception {
	result := tx.Save(&data)
	if result.Error != nil {
		return exceptions.ParseGormError(result.Error)
	}

	return nil
}

func (r *CompRepositoriesImpl) Delete(ctx *gin.Context, tx *gorm.DB, data models.Agents) *exceptions.Exception {
	result := tx.Delete(&data)
	if result.Error != nil {
		return exceptions.ParseGormError(result.Error)
	}

	return nil
}
