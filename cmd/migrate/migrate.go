package main

import (
	"nganterin-cs/pkg/config"
	"nganterin-cs/models"
)

func main() {
	db := config.InitDB()

	err := db.AutoMigrate(
		&models.Agents{},
	)
	if err != nil {
		panic("failed to migrate models: " + err.Error())
	}
}
