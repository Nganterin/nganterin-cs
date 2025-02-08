package controllers

import (
	"net/http"
	"nganterin-cs/internal/auth/dto"
	"nganterin-cs/internal/auth/services"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

type CompControllersImpl struct {
	services services.CompServices
}

func NewCompController(compServices services.CompServices) CompControllers {
	return &CompControllersImpl{
		services: compServices,
	}
}

func (h *CompControllersImpl) Login(ctx *gin.Context) {
	var data dto.Login

	errRequest := ctx.ShouldBindJSON(&data)
	if errRequest != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.NewException(http.StatusBadRequest, errRequest.Error()))
		return
	}

	token, err := h.services.Login(ctx, data)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Body:    token,
		Message: "login success",
	})
}
