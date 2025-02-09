package services

import (
	"log"
	"net/http"
	"nganterin-cs/api/agents/dto"
	"nganterin-cs/api/agents/repositories"
	"nganterin-cs/pkg/exceptions"
	"nganterin-cs/pkg/helpers"
	"nganterin-cs/pkg/mapper"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"

	emailDTO "nganterin-cs/emails/dto"
	emails "nganterin-cs/emails/services"
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

	go func() {
		err := emails.SendAgentAccountEmail(emailDTO.EmailAgentAccount{
			Email:    data.Email,
			Username: data.Username,
			Password: password,
		})
		if err != nil {
			log.Println("failed to send agent account email: " + err.Error())
		}
	}()

	return nil
}

func (s *CompServicesImpl) SignIn(ctx *gin.Context, data dto.Login) (*string, *exceptions.Exception) {
	validateErr := s.validate.Struct(data)
	if validateErr != nil {
		return nil, exceptions.NewValidationException(validateErr)
	}

	agentData, err := s.repo.FindByUsername(ctx, s.DB, data.Username)
	if err != nil {
		return nil, err
	}

	err = helpers.CheckPasswordHash(data.Password, agentData.HashedPassword)
	if err != nil {
		return nil, err
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["uuid"] = agentData.UUID
	claims["email"] = agentData.Email
	claims["name"] = agentData.Name
	claims["role"] = agentData.Role
	claims["type"] = "agent"

	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	secretKey := []byte(secret)
	tokenString, signError := token.SignedString(secretKey)
	if signError != nil {
		return nil, exceptions.NewException(http.StatusInternalServerError, exceptions.ErrTokenGenerate)
	}

	return &tokenString, nil
}
