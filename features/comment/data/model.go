package data

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	PostID    string `gorm:"not null"`
	UserID    string `gorm:"not null"`
	Content   string `gorm:"not null"`
	CreatedAt time.Time
}
