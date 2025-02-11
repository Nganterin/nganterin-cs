package models

import (
	"time"

	"gorm.io/gorm"
)

type Chats struct {
	gorm.Model

	ID           uint   `gorm:"primaryKey"`
	UUID         string `gorm:"not null;unique;index"`
	CustomerUUID string `gorm:"not null;index"`
	Message      string `gorm:"not null"`
	IsCSChat     bool   `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"null;default:null"`

	Customer Customers `gorm:"foreignKey:CustomerUUID;references:UUID"`
}
