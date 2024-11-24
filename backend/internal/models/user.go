package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	LastLogin    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
