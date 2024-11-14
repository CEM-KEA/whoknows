package models

import "time"

type SearchLog struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null"`
    Query     string    `gorm:"type:text;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}