package helpers

import (
	"net/http"
	"nganterin-cs/api/customers/dto"
	"nganterin-cs/pkg/exceptions"

	"github.com/gin-gonic/gin"
)

func GetCustomerData(c *gin.Context) (dto.CustomerOutput, *exceptions.Exception) {
	var result dto.CustomerOutput
	data, _ := c.Get("customer")

	result, ok := data.(dto.CustomerOutput)
	if !ok {
		return result, exceptions.NewException(http.StatusUnauthorized, exceptions.ErrInvalidTokenStructure)
	}

	return result, nil
}
