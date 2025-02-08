package mapper

import (
	"nganterin-cs/api/example/dto"
	"nganterin-cs/models"

	"github.com/go-viper/mapstructure/v2"
)

func MapExampleInputToModel(input dto.ExampleInput) models.Example {
	var example models.Example

	mapstructure.Decode(input, &example)
	return example
}
