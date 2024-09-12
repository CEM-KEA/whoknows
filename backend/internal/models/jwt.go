package models

import "time"

type JWT struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index;not null"`
	Token     string `gorm:"type:text;not null"`
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}
