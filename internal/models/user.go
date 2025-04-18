package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"required"`
	Email     string         `json:"email" gorm:"unique" binding:"required,email"`
	Password  string         `json:"-"`
	Role      string         `json:"role" gorm:"default:user"`
	Groups    []Group        `json:"groups" gorm:"many2many:group_users"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
