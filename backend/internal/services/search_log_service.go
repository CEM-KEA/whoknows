package services

import (
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"gorm.io/gorm"
)

// CreateSearchLog creates a new search log entry in the database.
// It logs any errors encountered during the creation process.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance used to interact with the database.
//   - searchLog: A pointer to the models.SearchLog instance containing the search log details.
//
// Returns:
//   - error: An error object if the creation fails, otherwise nil.
func CreateSearchLog(db *gorm.DB, searchLog *models.SearchLog) error {
	err := db.Create(searchLog).Error
	if err != nil {
		utils.LogError(err, "Failed to create search log", utils.SanitizeFields(map[string]interface{}{
			"query": searchLog.Query,
		}))
		return err
	}
	return nil
}
