package models

import (
	"gorm.io/gorm"
	"time"
)

type Group struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"unique;not null"`
	Users     []User         `json:"users" gorm:"many2many:group_users"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
