package services

import (
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"gorm.io/gorm"
)

// CreateSearchLog creates a new search log in the database.
func CreateSearchLog(db *gorm.DB, searchLog *models.SearchLog) error {
	return db.Create(searchLog).Error
}
