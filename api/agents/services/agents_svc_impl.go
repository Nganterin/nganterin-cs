package services

import (
	"nganterin-cs/api/agents/dto"
	"nganterin-cs/api/agents/repositories"
	"nganterin-cs/pkg/exceptions"
	"nganterin-cs/pkg/helpers"
	"nganterin-cs/pkg/mapper"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func (s *CompServicesImpl) Create(ctx *gin.Context, data dto.Agents) *exceptions.Exception {
	validateErr := s.validate.Struct(data)
	if validateErr != nil {
		return exceptions.NewValidationException(validateErr)
	}

	password := helpers.GenerateUniquePassword()
	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		return err
	}

	input := mapper.MapAgentInputToModel(data)
	input.UUID = uuid.NewString()
	input.HashedPassword = hashedPassword

	err = s.repo.Create(ctx, s.DB, input)
	if err != nil {
		return err
	}

	return nil
}
