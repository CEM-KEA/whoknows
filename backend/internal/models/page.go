package models

import (
	"time"
)

type Page struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"type:varchar(255);not null"`
	Url       string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Language  string `gorm:"type:varchar(2);not null;CHECK (language IN ('en', 'da')) DEFAULT 'en'"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
