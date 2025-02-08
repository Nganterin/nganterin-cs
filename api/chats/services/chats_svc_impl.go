package services

import (
	"nganterin-cs/api/chats/dto"
	"nganterin-cs/api/chats/repositories"
	"nganterin-cs/pkg/exceptions"
	"nganterin-cs/pkg/mapper"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CompServicesImpl struct {
	repo     repositories.CompRepositories
	DB       *gorm.DB
	validate *validator.Validate
}

func NewComponentServices(compRepositories repositories.CompRepositories, db *gorm.DB, validate *validator.Validate) CompServices {
	return &CompServicesImpl{
		repo:     compRepositories,
		DB:       db,
		validate: validate,
	}
}

func (s *CompServicesImpl) Create(ctx *gin.Context, data dto.Chats) *exceptions.Exception {
	validateErr := s.validate.Struct(data)
	if validateErr != nil {
		return exceptions.NewValidationException(validateErr)
	}

	input := mapper.MapChatInputToModel(data)

	err := s.repo.Create(ctx, s.DB, input)
	if err != nil {
		return err
	}

	return nil
}