package database

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Name  string `gorm:"column:name;primaryKey"`
	Color int    `gorm:"column:color;default:0"`
	Mood  string `gorm:"column:mood;default:😀"`
}

type Session struct {
	ID        uuid.UUID  `gorm:"column:id;primaryKey"`
	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime"`
	IsActive  bool       `gorm:"column:is_active;default:true"`
	UserName  string     `gorm:"column:user_name"`
	User      User
}
