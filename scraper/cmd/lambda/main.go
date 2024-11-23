package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/CEM-KEA/whoknows/backend/scraper/internal/db"
	"github.com/CEM-KEA/whoknows/backend/scraper/internal/wiki"
	"github.com/aws/aws-lambda-go/lambda"
)

// getScraperConfig initializes and returns a pointer to a wiki.ScraperConfig struct.
// It reads configuration values from environment variables and provides default values
// if the environment variables are not set.
//
// Environment Variables:
// - SCRAPER_MAX_REDIRECTS: Maximum number of redirects allowed (default: 3)
// - SCRAPER_TIMEOUT: Timeout duration for the scraper (default: 30s)
// - SCRAPER_RATE_LIMIT: Rate limit duration between requests (default: 1s)
// - SCRAPER_MAX_DEPTH: Maximum depth for scraping (default: 1)
// - SCRAPER_PARALLEL_REQUESTS: Number of parallel requests allowed (default: 2)
//
// Returns:
// - *wiki.ScraperConfig: A pointer to the initialized ScraperConfig struct.
func getScraperConfig() *wiki.ScraperConfig {
	parseEnvInt := func(key string, defaultValue int) int {
		if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
			return val
		}
		return defaultValue
	}

	parseEnvDuration := func(key string, defaultValue time.Duration) time.Duration {
		if val, err := time.ParseDuration(os.Getenv(key)); err == nil {
			return val
		}
		return defaultValue
	}

	return &wiki.ScraperConfig{
		MaxRedirects:     parseEnvInt("SCRAPER_MAX_REDIRECTS", 3),
		TimeoutDuration:  parseEnvDuration("SCRAPER_TIMEOUT", 30*time.Second),
		RateLimit:        parseEnvDuration("SCRAPER_RATE_LIMIT", 1*time.Second),
		MaxDepth:         parseEnvInt("SCRAPER_MAX_DEPTH", 1),
		ParallelRequests: parseEnvInt("SCRAPER_PARALLEL_REQUESTS", 2),
	}
}

// handleRequest is the main entry point for processing the scraping task.
// It initializes metrics, connects to the database, fetches queries, and processes each row using the scraper.
// The function logs metrics and handles errors appropriately.
//
// Parameters:
// - ctx: The context for managing request-scoped values, cancellation, and deadlines.
//
// Returns:
// - error: An error if any step in the process fails, otherwise nil.
func handleRequest(ctx context.Context) error {
    metrics := wiki.NewMetrics()
    defer metrics.Log()

    database, err := db.Connect(ctx)
    if err != nil {
        return err
    }
    defer database.Close()

    rows, err := db.FetchQueries(ctx, database)
    if err != nil {
        return err
    }
    defer rows.Close()

    scraper := wiki.NewScraper(getScraperConfig())
    for rows.Next() {
        metrics.QueriesProcessed++
        if err := scraper.ProcessRow(ctx, database, rows); err != nil {
            metrics.FailedScrapes++
            log.Printf("Failed to process row: %v", err)
            continue
        }
        metrics.SuccessfulScrapes++
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating through rows: %v", err)
        return err
    }

    return nil
}


// main is the entry point for the AWS Lambda function. It starts the Lambda
// function by calling lambda.Start with the handleRequest function as the
// handler.
func main() {
    lambda.Start(handleRequest)
}