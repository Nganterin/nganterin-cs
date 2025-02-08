package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	General Role = "general"
)

type Agents struct {
	gorm.Model

	ID             uint   `gorm:"primaryKey"`
	UUID           string `gorm:"not null;unique;index"`
	Username       string `gorm:"not null;unique;index"`
	HashedPassword string `gorm:"not null"`
	Email          string `gorm:"not null;unique;index"`
	Role           Role   `gorm:"type:agent_role;not null"`
	IsActive       bool   `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"null;default:null"`
}
