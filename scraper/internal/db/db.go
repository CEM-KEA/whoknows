package db

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"
    _ "github.com/lib/pq"
)

// Connect establishes a connection to the PostgreSQL database using the
// connection details provided via environment variables. It returns a
// sql.DB object representing the connection pool and an error if the
// connection could not be established.
//
// Environment Variables:
// - API_DATABASE_HOST: The hostname of the database server.
// - API_DATABASE_PORT: The port number on which the database server is listening.
// - API_DATABASE_USER: The username to use for authentication.
// - API_DATABASE_PASSWORD: The password to use for authentication.
// - API_DATABASE_NAME: The name of the database to connect to.
//
// The function configures the connection pool with a maximum connection
// lifetime of 3 minutes, a maximum of 10 open connections, and a maximum
// of 10 idle connections. It also attempts to ping the database with
// retries to ensure the connection is valid.
//
// Parameters:
// - ctx: A context.Context object for managing the connection lifecycle.
//
// Returns:
// - *sql.DB: A pointer to the sql.DB object representing the connection pool.
// - error: An error object if the connection could not be established.
func Connect(ctx context.Context) (*sql.DB, error) {
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

    if err := pingWithRetry(ctx, db); err != nil {
        return nil, err
    }

    return db, nil
}

// FetchQueries retrieves a list of search queries from the database that have not been scraped
// in the last 48 hours or have never been scraped. The queries are ranked by the number of times
// they appear in the search logs, and the top 20 queries are returned.
//
// Parameters:
//   ctx - The context to use for the database query.
//   db  - The database connection to use.
//
// Returns:
//   *sql.Rows - The result set containing the id, query, and count of the top 20 queries.
//   error     - An error if the query fails.
func FetchQueries(ctx context.Context, db *sql.DB) (*sql.Rows, error) {
    return db.QueryContext(ctx, `
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
}

// StorePage stores a page in the database and updates the associated queries.
// It starts a transaction, stores the page, updates the queries, and commits the transaction.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - db: The database connection.
//   - page: The page to be stored.
//   - queryID: The ID of the query to be updated.
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
func StorePage(ctx context.Context, db *sql.DB, page *Page, queryID uint) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("error starting transaction: %v", err)
    }
    defer tx.Rollback()

    if err := storePage(ctx, tx, page); err != nil {
        return err
    }

    if err := updateQueries(ctx, tx, queryID); err != nil {
        return err
    }

    return tx.Commit()
}

// storePage inserts a new page into the 'pages' table or updates the existing page if a conflict on the URL occurs.
// It uses a transaction to ensure atomicity and consistency.
//
// Parameters:
//  - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//  - tx: The transaction object for executing the SQL statement.
//  - page: A pointer to the Page struct containing the page details to be stored.
//
// Returns:
//  - error: An error object if an error occurs during the execution of the SQL statement, otherwise nil.
func storePage(ctx context.Context, tx *sql.Tx, page *Page) error {
    _, err := tx.ExecContext(ctx, `
        INSERT INTO pages (title, url, language, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (url) DO UPDATE 
        SET content = EXCLUDED.content,
            updated_at = EXCLUDED.updated_at
    `, page.Title, page.URL, page.Language, page.Content, page.CreatedAt, page.UpdatedAt)
    return err
}

// updateQueries updates the 'scraped_at' timestamp for rows in the 'search_logs' table
// that match the query text associated with the given queryID. The query text is 
// retrieved from the 'search_logs' table using the provided queryID. The function 
// logs the number of rows affected by the update.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation, and timeouts.
//   - tx: The transaction within which the update operation is performed.
//   - queryID: The ID of the query whose text is used to find matching rows to update.
//
// Returns:
//   - error: An error object if any error occurs during the execution of the function, 
//            otherwise nil.
func updateQueries(ctx context.Context, tx *sql.Tx, queryID uint) error {
    var queryText string
    err := tx.QueryRowContext(ctx, `
        SELECT query FROM search_logs WHERE id = $1
    `, queryID).Scan(&queryText)
    if err != nil {
        return err
    }

    result, err := tx.ExecContext(ctx, `
        UPDATE search_logs 
        SET scraped_at = NOW()
        WHERE LOWER(TRIM(query)) = LOWER(TRIM($1))
    `, queryText)
    if err != nil {
        return err
    }

    if count, err := result.RowsAffected(); err == nil {
        log.Printf("Updated %d queries matching: %s", count, queryText)
    }
    return nil
}

// pingWithRetry attempts to ping the database with a specified number of retries.
// It takes a context and a database connection as parameters.
// If the ping is successful within the retry limit, it returns nil.
// If all retries are exhausted, it returns an error indicating the failure.
//
// Parameters:
//   ctx - the context to control the timeout and cancellation of the ping operation
//   db - the database connection to be pinged
//
// Returns:
//   error - nil if the ping is successful, or an error if the ping fails after the maximum retries
func pingWithRetry(ctx context.Context, db *sql.DB) error {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        if err := db.PingContext(ctx); err == nil {
            return nil
        }
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return fmt.Errorf("database ping failed after %d retries", maxRetries)
}