package models

import "time"

type SearchLog struct {
    ID        uint      `gorm:"primaryKey"`
    Query     string    `gorm:"type:text;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}