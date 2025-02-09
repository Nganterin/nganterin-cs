package services

import (
	"net/http"
	"nganterin-cs/api/customers/dto"
	"nganterin-cs/api/customers/repositories"
	"nganterin-cs/pkg/exceptions"
	"nganterin-cs/pkg/mapper"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func (s *CompServicesImpl) Create(ctx *gin.Context, data dto.Customers) (*string, *exceptions.Exception) {
	validateErr := s.validate.Struct(data)
	if validateErr != nil {
		return nil, exceptions.NewValidationException(validateErr)
	}

	input := mapper.MapCustomerInputToModel(data)
	input.UUID = uuid.NewString()

	err := s.repo.Create(ctx, s.DB, input)
	if err != nil {
		return nil, err
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["uuid"] = input.UUID
	claims["email"] = data.Email
	claims["name"] = data.Name
	claims["phone"] = data.Phone
	claims["type"] = "customer"

	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	secretKey := []byte(secret)
	tokenString, signError := token.SignedString(secretKey)
	if signError != nil {
		return nil, exceptions.NewException(http.StatusInternalServerError, exceptions.ErrTokenGenerate)
	}

	return &tokenString, nil
}
