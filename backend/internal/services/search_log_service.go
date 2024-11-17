package services

import (
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateSearchLog creates a new search log entry in the database.
// It logs the creation process and any errors that occur.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance used to interact with the database.
//   - searchLog: A pointer to the models.SearchLog instance containing the search log details.
//
// Returns:
//   - error: An error object if the creation fails, otherwise nil.
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
