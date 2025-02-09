package mapper

import (
	"nganterin-cs/api/customers/dto"
	"nganterin-cs/models"

	"github.com/go-viper/mapstructure/v2"
)

func MapCustomerInputToModel(data dto.Customers) models.Customers {
	var result models.Customers

	mapstructure.Decode(data, &result)
	return result
}

func MapCustomerModelToOutput(data models.Customers) dto.CustomerOutput {
	var result dto.CustomerOutput

	mapstructure.Decode(data, &result)
	return result
}
