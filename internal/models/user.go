package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"required"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string         `json:"-"`
	RoleID       uint           `json:"role_id"`
	Role         Role           `json:"role" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsBanned     bool           `json:"is_banned" gorm:"default:false"`
	Groups       []Group        `json:"groups" gorm:"many2many:group_users"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
