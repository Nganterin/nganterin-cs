package main

import (
	"nganterin-cs/models"
	"nganterin-cs/pkg/config"
)

func main() {
	db := config.InitDB()

	if err := db.Exec(`DO $$ 
    	BEGIN
    	    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'agent_role') THEN
    	        CREATE TYPE agent_role AS ENUM ('general');
    	    ELSE
    	        ALTER TYPE agent_role ADD VALUE IF NOT EXISTS 'general';
    	    END IF;
    	END $$;`).Error; err != nil {
		panic("failed to create or update enum type: " + err.Error())
	}

	err := db.AutoMigrate(
		&models.Agents{},
		&models.Customers{},
	)
	if err != nil {
		panic("failed to migrate models: " + err.Error())
	}
}
