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

func (r *CompRepositoriesImpl) Create(ctx *gin.Context, tx *gorm.DB, data models.Chats) *exceptions.Exception {
	result := tx.
		Create(&data)
	if result.Error != nil {
		return exceptions.ParseGormError(tx, result.Error)
	}

	return nil
}

func (r *CompRepositoriesImpl) FindByUUID(ctx *gin.Context, tx *gorm.DB, uuid string) (*models.Chats, *exceptions.Exception) {
	var chat models.Chats

	result := tx.
		Preload("Customer").
		Where("uuid = ?", uuid).
		First(&chat)
	if result.Error != nil {
		return nil, exceptions.ParseGormError(tx, result.Error)
	}

	return &chat, nil
}

func (r *CompRepositoriesImpl) FindAll(ctx *gin.Context, tx *gorm.DB) ([]models.Chats, *exceptions.Exception) {
	var chats []models.Chats

	result := tx.
		Preload("Customer").
		Order("id ASC").
		Find(&chats)
	if result.Error != nil {
		return nil, exceptions.ParseGormError(tx, result.Error)
	}

	return chats, nil
}

func (r *CompRepositoriesImpl) FindAllByLastUUID(ctx *gin.Context, tx *gorm.DB, uuid string) ([]models.Chats, *exceptions.Exception) {
	var chats []models.Chats

	lastChat, err := r.FindByUUID(ctx, tx, uuid)
	if err != nil {
		return nil, err
	}

	result := tx.
		Preload("Customer").
		Where("id > ?", lastChat.ID).
		Order("id ASC").
		Find(&chats)
	if result.Error != nil {
		return nil, exceptions.ParseGormError(tx, result.Error)
	}

	return chats, nil
}

func (r *CompRepositoriesImpl) FindAllByCustomerUUID(ctx *gin.Context, tx *gorm.DB, uuid string) ([]models.Chats, *exceptions.Exception) {
	var chats []models.Chats

	result := tx.
		Preload("Customer").
		Where("customer_uuid = ?", uuid).
		Order("id ASC").
		Find(&chats)
	if result.Error != nil {
		return nil, exceptions.ParseGormError(tx, result.Error)
	}

	return chats, nil
}

func (r *CompRepositoriesImpl) FindAllByLastUUIDAndCustomerUUID(ctx *gin.Context, tx *gorm.DB, lastUUID string, customerUUID string) ([]models.Chats, *exceptions.Exception) {
	var chats []models.Chats

	lastChat, err := r.FindByUUID(ctx, tx, lastUUID)
	if err != nil {
		return nil, err
	}

	result := tx.
		Preload("Customer").
		Where("id > ?", lastChat.ID).
		Where("customer_uuid = ?", customerUUID).
		Order("id ASC").
		Find(&chats)
	if result.Error != nil {
		return nil, exceptions.ParseGormError(tx, result.Error)
	}

	return chats, nil
}
