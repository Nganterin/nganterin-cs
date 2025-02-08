package controllers

import "github.com/gin-gonic/gin"

type CompControllers interface {
	ChatWebSocket(ctx *gin.Context)
}
