package models

import (
	"time"

	"gorm.io/gorm"
)

type Customers struct {
	gorm.Model

	ID       uint   `gorm:"primaryKey"`
	UUID     string `gorm:"not null;unique;index"`
	Email    string `gorm:"not null;index"`
	Name     string `gorm:"not null"`
	Phone    string `gorm:"not null"`
	IsActive bool   `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"null;default:null"`
}
