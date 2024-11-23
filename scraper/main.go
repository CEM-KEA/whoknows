package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gocolly/colly/v2"
	_ "github.com/lib/pq"
)

type Page struct {
	ID        uint
	Title     string
	URL       string
	Language  string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ScraperMetrics struct {
    QueriesProcessed int
    SuccessfulScrapes int
    FailedScrapes int
    StartTime time.Time
    EndTime time.Time
}
// logMetrics logs the metrics of a scraper operation.
// It logs the duration of the scraping process, the number of queries processed,
// the number of successful scrapes, and the number of failed scrapes.
//
// Parameters:
//   - metrics: ScraperMetrics containing the start and end time of the scraping
//              process, the number of queries processed, the number of successful
//              scrapes, and the number of failed scrapes.
//   - error: an error if the query is invalid, otherwise nil
func logMetrics(metrics ScraperMetrics) {
    duration := metrics.EndTime.Sub(metrics.StartTime)
    log.Printf("Scraper finished in %v", duration)
    log.Printf("Queries processed: %d", metrics.QueriesProcessed)
    log.Printf("Successful scrapes: %d", metrics.SuccessfulScrapes)
    log.Printf("Failed scrapes: %d", metrics.FailedScrapes)
}

// validateQuery checks if the provided query string is valid.
// A valid query is non-empty and does not contain any of the following characters:
// <, >, {, }, |, \, ^, [, ], `
// If the query is empty or contains any invalid characters, an error is returned.
//
// Parameters:
//   - query: the query string to validate
//
// Returns:
//   - error: an error if the query is invalid, otherwise nil
func validateQuery(query string) error {
    if strings.TrimSpace(query) == "" {
        return fmt.Errorf("empty query")
    }
	
    invalidChars := []string{"<", ">", "{", "}", "|", "\\", "^", "[", "]", "`"}
    for _, char := range invalidChars {
        if strings.Contains(query, char) {
            return fmt.Errorf("query contains invalid character: %s", char)
        }
    }
    
    return nil
}

// handleRequest processes database queries using a web scraper and logs metrics.
// It connects to the database, fetches queries, and processes each row using a collector.
// Metrics such as start time, end time, number of queries processed, successful scrapes,
// and failed scrapes are recorded and logged.
//
// Parameters:
//   - ctx: The context for managing request deadlines, cancelation signals, and other request-scoped values.
//
// Returns:
//   - error: An error if any step in the process fails, otherwise nil.
func handleRequest(ctx context.Context) error {
	metrics := ScraperMetrics{
        StartTime: time.Now(),
    }
    defer func() {
        metrics.EndTime = time.Now()
        logMetrics(metrics)
    }()

	db, err := connectToDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := fetchQueries(ctx, db)
	if err != nil {
		return err
	}
	defer rows.Close()

	c := setupCollector()

	for rows.Next() {
        metrics.QueriesProcessed++
        if err := processRow(ctx, db, c, rows); err != nil {
            metrics.FailedScrapes++
            log.Printf("Error processing row: %v", err)
            continue
        }
        metrics.SuccessfulScrapes++
    }

	c.Wait()
	return nil
}

// connectToDB establishes a connection to a PostgreSQL database using the provided context.
// It retrieves the database connection details from environment variables and configures
// the connection pool settings. The function attempts to ping the database up to three times
// with increasing delays between retries if the initial ping fails.
//
// Parameters:
//   - ctx: A context to control the lifetime of the database connection attempt.
//
// Returns:
//   - *sql.DB: A pointer to the established database connection.
//   - error: An error if the connection or ping attempts fail.
func connectToDB(ctx context.Context) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("API_DATABASE_HOST"),
        os.Getenv("API_DATABASE_PORT"),
        os.Getenv("API_DATABASE_USER"),
        os.Getenv("API_DATABASE_PASSWORD"),
        os.Getenv("API_DATABASE_NAME"),
    )

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("database connection error: %v", err)
    }

    db.SetConnMaxLifetime(time.Minute * 3)
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(10)

    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        if err = db.PingContext(ctx); err == nil {
            break
        }
        time.Sleep(time.Second * time.Duration(i+1))
        log.Printf("Retry %d: Database ping failed: %v", i+1, err)
    }
    if err != nil {
        return nil, fmt.Errorf("database ping error after %d retries: %v", maxRetries, err)
    }

    return db, nil
}

// fetchQueries retrieves a set of queries from the database that need to be scraped.
// It selects queries from the `search_logs` table that have either never been scraped
// or were last scraped more than 48 hours ago. The queries are ranked by the number
// of times they appear in the log, and the top 20 queries are returned.
//
// Parameters:
//  - ctx: The context for managing request deadlines and cancellation signals.
//  - db: The database connection.
//
// Returns:
//  - *sql.Rows: A result set containing the selected queries.
//  - error: An error object if the query execution fails.
func fetchQueries(ctx context.Context, db *sql.DB) (*sql.Rows, error) {
	rows, err := db.QueryContext(ctx, `
		WITH ranked_queries AS (
			SELECT 
				id,
				query,
				COUNT(*) as count,
				MAX(scraped_at) as last_scraped
			FROM search_logs
			GROUP BY id, query
			HAVING MAX(scraped_at) IS NULL 
				OR MAX(scraped_at) < NOW() - INTERVAL '48 hours'
			ORDER BY count DESC
			LIMIT 20
		)
		SELECT id, query, count
		FROM ranked_queries
	`)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	return rows, nil
}

// setupCollector initializes and configures a new Colly collector with specific settings.
// It sets the maximum depth for crawling, user agent, allows URL revisits, enables asynchronous requests,
// and restricts the allowed domains to Wikipedia. Additionally, it sets a limit rule for the domain
// with a parallelism of 2 and a delay of 1 second between requests.
//
// Returns:
//   *colly.Collector: A pointer to the configured Colly collector.
func setupCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.AllowedDomains("wikipedia.org", "www.wikipedia.org", "en.wikipedia.org"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*wikipedia.org",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	return c
}

// processRow processes a single row from the database, validates the query, scrapes the corresponding Wikipedia page,
// and stores the scraped content in the database.
//
// Parameters:
// - ctx: The context for controlling the lifecycle of the scraping process.
// - db: The database connection to use for storing the scraped content and updating the query log.
// - c: The Colly collector used for scraping the web page.
// - rows: The SQL rows containing the query information to process.
//
// Returns:
// - error: An error if any step of the process fails, otherwise nil.
func processRow(ctx context.Context, db *sql.DB, c *colly.Collector, rows *sql.Rows) error {
    var queryID uint
    var query string
    var count int

    if err := rows.Scan(&queryID, &query, &count); err != nil {
        return fmt.Errorf("error scanning row: %v", err)
    }

    if err := validateQuery(query); err != nil {
        return fmt.Errorf("invalid query %s: %v", query, err)
    }

    page := &Page{
        Language:  "en",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    c.OnHTML("article#content", func(e *colly.HTMLElement) {
        page.Title = strings.TrimSpace(e.ChildText("h1#firstHeading"))
        if page.Title == "" {
            page.Title = query
        }

        content := e.ChildText("div#mw-content-text")
        content = strings.ReplaceAll(content, "[edit]", "")
        content = strings.ReplaceAll(content, "[citation needed]", "")
        page.Content = strings.TrimSpace(content)
    })

    c.OnError(func(r *colly.Response, err error) {
        log.Printf("Error scraping %s: %v", r.Request.URL, err)
    })

    targetURL := "https://en.wikipedia.org/wiki/" + url.QueryEscape(query)
    page.URL = targetURL

    scrapeCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    if err := c.Visit(targetURL); err != nil {
        return fmt.Errorf("error visiting %s: %v", targetURL, err)
    }

    return storePageAndUpdateQueryLog(scrapeCtx, db, page, queryID)
}


// storePageAndUpdateQueryLog stores a page in the database and updates the query log with the current timestamp.
// It performs these operations within a single transaction to ensure atomicity.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - db: The database connection.
//   - page: The page to be stored, containing title, URL, language, content, created_at, and updated_at fields.
//   - queryID: The ID of the query log entry to be updated.
//
// Returns:
//   - error: An error object if any operation fails, otherwise nil.
func storePageAndUpdateQueryLog(ctx context.Context, db *sql.DB, page *Page, queryID uint) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO pages (title, url, language, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (url) DO UPDATE 
		SET content = EXCLUDED.content,
			updated_at = EXCLUDED.updated_at
	`, page.Title, page.URL, page.Language, page.Content, page.CreatedAt, page.UpdatedAt)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error storing page: %v", err)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE query_log 
		SET scraped_at = NOW()
		WHERE id = $1
	`, queryID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating scraped_at: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// main is the entry point of the application. It starts the AWS Lambda function
// by calling lambda.Start with the handleRequest function as the handler.
func main() {
	lambda.Start(handleRequest)
}