package database

import (
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"gorm.io/gorm"
)

// RegisterCallbacks registers GORM callbacks for Prometheus metrics collection.
// It sets up the following callbacks:
// - After querying the database, the "prometheus:after_query" callback is triggered.
// - After creating a record in the database, the "prometheus:after_create" callback is triggered.
// - After updating a record in the database, the "prometheus:after_update" callback is triggered.
// - After deleting a record from the database, the "prometheus:after_delete" callback is triggered.
//
// These callbacks are used to collect and report metrics to Prometheus.
func RegisterCallbacks(db *gorm.DB) {
	db.Callback().Query().After("gorm:query").Register("prometheus:after_query", afterQuery)
	db.Callback().Create().After("gorm:create").Register("prometheus:after_create", afterCreate)
	db.Callback().Update().After("gorm:update").Register("prometheus:after_update", afterUpdate)
	db.Callback().Delete().After("gorm:delete").Register("prometheus:after_delete", afterDelete)
}

// Callback function to observe query duration and failures
func afterQuery(db *gorm.DB) {
	observeGormOperation("query", db)
}

func afterCreate(db *gorm.DB) {
	observeGormOperation("create", db)
}

func afterUpdate(db *gorm.DB) {
	observeGormOperation("update", db)
}

func afterDelete(db *gorm.DB) {
	observeGormOperation("delete", db)
}

// observeGormOperation observes the duration of a GORM database operation and records metrics.
// It takes the operation name as a string and a GORM database instance as parameters.
// The function measures the time taken for the operation to complete and records the duration
// using the ObserveDBQueryDuration utility function. If the operation results in an error,
// it increments the failed queries counter using the DBFailedQueries utility.
func observeGormOperation(operation string, db *gorm.DB) {
	start := time.Now()

	// Defer to calculate duration after the operation
	defer func() {
		duration := time.Since(start).Seconds()
		utils.ObserveDBQueryDuration(operation, duration)
		if db.Error != nil {
			utils.DBFailedQueries.WithLabelValues(operation).Inc()
		}
	}()
}
