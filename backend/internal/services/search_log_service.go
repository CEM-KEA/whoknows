package services

import (
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateSearchLog creates a new search log in the database.
func CreateSearchLog(db *gorm.DB, searchLog *models.SearchLog) error {
	utils.LogInfo("Creating search log", logrus.Fields{
		"query": searchLog.Query,
	})

	err := db.Create(searchLog).Error
	if err != nil {
		utils.LogError(err, "Failed to create search log", logrus.Fields{
			"query": searchLog.Query,
		})
		return err
	}

	utils.LogInfo("Search log created successfully", logrus.Fields{
		"query": searchLog.Query,
		"id":    searchLog.ID,
	})
	return nil
}
