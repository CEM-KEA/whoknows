package wiki

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/CEM-KEA/whoknows/backend/scraper/internal/db"
	"github.com/gocolly/colly/v2"
)

type Scraper struct {
	collector *colly.Collector
	config    *ScraperConfig
}

type ScraperConfig struct {
	MaxRedirects     int
	TimeoutDuration  time.Duration
	RateLimit        time.Duration
	MaxDepth         int
	ParallelRequests int
}

type pageContent struct {
	title   string
	content string
}

// NewScraper creates a new Scraper instance with the provided configuration.
// If the provided configuration is nil, it initializes a default configuration.
//
// Parameters:
//   - config: A pointer to a ScraperConfig struct containing the scraper settings.
//
// Returns:
//   - A pointer to a Scraper struct initialized with the provided or default configuration.
func NewScraper(config *ScraperConfig) *Scraper {
	if config == nil {
		config = &ScraperConfig{
			MaxRedirects:     3,
			TimeoutDuration:  30 * time.Second,
			RateLimit:        1 * time.Second,
			MaxDepth:         1,
			ParallelRequests: 2,
		}
	}

	c := colly.NewCollector(
		colly.MaxDepth(config.MaxDepth),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.AllowedDomains("wikipedia.org", "www.wikipedia.org", "en.wikipedia.org"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*wikipedia.org",
		Parallelism: config.ParallelRequests,
		Delay:       config.RateLimit,
	})

	return &Scraper{
		collector: c,
		config:    config,
	}
}

// ProcessRow processes a single row from the database, performing web scraping
// based on the query contained in the row. It scans the row to extract the query,
// preprocesses the query, and then uses a web scraper to search for relevant pages.
// If a valid page is found, it stores the page in the database.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - database: The database connection to store the scraped page.
//   - rows: The database rows containing the query to be processed.
//
// Returns:
//   - error: An error if any issues occur during processing, or nil if successful.
func (s *Scraper) ProcessRow(ctx context.Context, database *sql.DB, rows *sql.Rows) error {
	queryID, query, err := s.scanRow(rows)
	if err != nil {
		return err
	}

	processedQuery, valid := preprocessQuery(query)
	if !valid {
		log.Printf("Skipping invalid query: %s", query)
		return nil
	}

	page := db.NewPage()
	foundValidPage := false
	redirectCount := 0
	visitedURLs := make(map[string]bool)

	s.configureCollectorCallbacks(query, page, &foundValidPage, &redirectCount, visitedURLs)

	timeoutCtx, cancel := context.WithTimeout(ctx, s.config.TimeoutDuration)
	defer cancel()

	urls := buildSearchURLs(processedQuery, query)
	for _, targetURL := range urls {
		if foundValidPage {
			break
		}

		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout while processing query: %s", query)
		default:
			if err := s.collector.Visit(targetURL); err != nil {
				log.Printf("Error visiting %s: %v", targetURL, err)
				continue
			}
			time.Sleep(s.config.RateLimit) // Use configured rate limit
		}
	}

	if foundValidPage {
		return db.StorePage(ctx, database, page, queryID)
	}

	log.Printf("No valid content found for query: %s", query)
	return nil
}

// scanRow scans a single row from the provided sql.Rows object and extracts
// the query ID, query string, and count. It returns the query ID and query
// string, or an error if the row scanning fails.
//
// Parameters:
//   rows - a pointer to an sql.Rows object containing the row to be scanned.
//
// Returns:
//   uint - the query ID extracted from the row.
//   string - the query string extracted from the row.
//   error - an error if the row scanning fails, otherwise nil.
func (s *Scraper) scanRow(rows *sql.Rows) (uint, string, error) {
	var queryID uint
	var query string
	var count int

	if err := rows.Scan(&queryID, &query, &count); err != nil {
		return 0, "", fmt.Errorf("error scanning row: %v", err)
	}

	return queryID, query, nil
}

// configureCollectorCallbacks sets up the callbacks for the collector to handle HTML elements.
// It processes the search results, handles redirects, and extracts content from the page.
//
// Parameters:
//   - query: The search query string.
//   - page: A pointer to the db.Page struct where the extracted content will be stored.
//   - foundValidPage: A pointer to a boolean that indicates if a valid page has been found.
//   - redirectCount: A pointer to an integer that tracks the number of redirects.
//   - visitedURLs: A map that keeps track of visited URLs to avoid processing the same URL multiple times.
func (s *Scraper) configureCollectorCallbacks(query string, page *db.Page, foundValidPage *bool, redirectCount *int, visitedURLs map[string]bool) {
	s.collector.OnHTML("body", func(e *colly.HTMLElement) {
		if *foundValidPage {
			return
		}

		if s.handleSearchResults(e, query, visitedURLs, redirectCount) {
			return
		}

		if s.handleRedirects(e, visitedURLs, redirectCount) {
			return
		}

		if content := s.tryExtractContent(e); content != nil {
			*foundValidPage = true
			page.Title = content.title
			page.Content = content.content
			page.URL = e.Request.URL.String()
		}
	})
}

// handleSearchResults processes the search results from a Wikipedia search query.
// It checks if there are any search results and finds the best match for the query.
// If a valid result is found and it hasn't been visited yet, it marks the URL as visited,
// increments the redirect count, and visits the URL if the maximum number of redirects
// has not been exceeded.
//
// Parameters:
//   - e: The HTMLElement containing the search results.
//   - query: The search query string.
//   - visitedURLs: A map to keep track of visited URLs to avoid revisiting.
//   - redirectCount: A pointer to an integer tracking the number of redirects.
//
// Returns:
//   - bool: Returns true if search results were found and processed, otherwise false.
func (s *Scraper) handleSearchResults(e *colly.HTMLElement, query string, visitedURLs map[string]bool, redirectCount *int) bool {
	if searchResults := e.ChildText(".mw-search-results"); searchResults != "" {
		bestResult := s.findBestSearchResult(e, query)
		if bestResult != "" && !visitedURLs[bestResult] {
			visitedURLs[bestResult] = true
			*redirectCount++
			if *redirectCount <= s.config.MaxRedirects {
				e.Request.Visit("https://en.wikipedia.org" + bestResult)
			}
		}
		return true
	}
	return false
}

// handleRedirects handles the redirection logic for the scraper.
// It checks if the current element contains a redirect link and if the link has not been visited before.
// If the redirect link is valid and within the maximum allowed redirects, it visits the redirect link.
//
// Parameters:
// - e: The current HTML element being processed.
// - visitedURLs: A map to keep track of visited URLs to avoid cycles.
// - redirectCount: A pointer to an integer that tracks the number of redirects.
//
// Returns:
// - bool: Returns true if a redirect was handled, otherwise false.
func (s *Scraper) handleRedirects(e *colly.HTMLElement, visitedURLs map[string]bool, redirectCount *int) bool {
	if redirTarget := e.ChildAttr(".redirectText a", "href"); redirTarget != "" {
		if !visitedURLs[redirTarget] {
			visitedURLs[redirTarget] = true
			*redirectCount++
			if *redirectCount <= s.config.MaxRedirects {
				e.Request.Visit("https://en.wikipedia.org" + redirTarget)
			}
		}
		return true
	}
	return false
}

// tryExtractContent attempts to extract the title and content from the provided HTML element.
// It trims any whitespace from the title and retrieves the content from the specified HTML tags.
// If the extracted title and content are not valid, it returns nil.
// Otherwise, it cleans the content and returns a pointer to a pageContent struct containing the title and cleaned content.
//
// Parameters:
//   e (*colly.HTMLElement): The HTML element from which to extract the content.
//
// Returns:
//   *pageContent: A pointer to a pageContent struct containing the extracted title and cleaned content,
//                 or nil if the content is not valid.
func (s *Scraper) tryExtractContent(e *colly.HTMLElement) *pageContent {
	title := strings.TrimSpace(e.ChildText("h1#firstHeading"))
	content := e.ChildText("div#mw-content-text")

	if !isValidContent(title, content) {
		return nil
	}

	content = cleanContent(content)
	return &pageContent{
		title:   title,
		content: content,
	}
}

// findBestSearchResult iterates over search results on a Wikipedia search result page
// and finds the best match for the given query. It calculates a match score for each
// result based on the similarity between the query terms and the result title. The
// function returns the link to the best matching search result.
//
// Parameters:
// - e: A colly.HTMLElement representing the search result page.
// - query: A string containing the search query.
//
// Returns:
// - A string containing the link to the best matching search result.
func (s *Scraper) findBestSearchResult(e *colly.HTMLElement, query string) string {
	var bestMatch string
	var bestScore float64

	queryTerms := strings.Fields(strings.ToLower(query))

	e.ForEach(".mw-search-result", func(_ int, s *colly.HTMLElement) {
		title := strings.ToLower(s.ChildText(".mw-search-result-heading"))
		link := s.ChildAttr(".mw-search-result-heading a", "href")

		score := calculateMatchScore(queryTerms, strings.Fields(title))
		if score > bestScore {
			bestScore = score
			bestMatch = link
		}
	})

	return bestMatch
}

// isValidContent checks if the provided title and content are valid.
// It returns false if the title or content is empty, or if the content length is less than 200 characters.
// Additionally, it returns false if the content contains any of the predefined invalid patterns.
// The function returns true if none of these conditions are met.
func isValidContent(title, content string) bool {
	if title == "" || content == "" || len(content) < 200 {
		return false
	}

	invalidPatterns := []string{
		"Wikipedia does not have an article",
		"Wikipedia does not yet have an article",
		"This article needs to be created",
		"If you expected a page to be here",
		"may refer to:", // Skip disambiguation pages
	}

	for _, pattern := range invalidPatterns {
		if strings.Contains(content, pattern) {
			return false
		}
	}

	return true
}

// cleanContent removes specific unwanted substrings from the given content string.
// It removes "[edit]", "[citation needed]", "Jump to navigation", and "Jump to search".
// Finally, it trims any leading and trailing whitespace from the cleaned content.
// 
// Parameters:
//   content (string): The input string to be cleaned.
//
// Returns:
//   string: The cleaned content string.
func cleanContent(content string) string {
	content = strings.ReplaceAll(content, "[edit]", "")
	content = strings.ReplaceAll(content, "[citation needed]", "")
	content = strings.ReplaceAll(content, "Jump to navigation", "")
	content = strings.ReplaceAll(content, "Jump to search", "")
	return strings.TrimSpace(content)
}

// calculateMatchScore calculates a match score between two sets of terms.
// It compares each term in queryTerms with each term in titleTerms and counts
// the number of matches where one term contains the other. The match score is
// then calculated as the ratio of the match count to the maximum length of the
// two term lists.
//
// Parameters:
//   - queryTerms: a slice of strings representing the query terms.
//   - titleTerms: a slice of strings representing the title terms.
//
// Returns:
//   - A float64 representing the match score, which is the ratio of the match
//     count to the maximum length of the two term lists. If either list is empty,
//     the function returns 0.
func calculateMatchScore(queryTerms, titleTerms []string) float64 {
	matchCount := 0
	for _, qTerm := range queryTerms {
		for _, tTerm := range titleTerms {
			if strings.Contains(tTerm, qTerm) || strings.Contains(qTerm, tTerm) {
				matchCount++
			}
		}
	}

	if len(queryTerms) == 0 || len(titleTerms) == 0 {
		return 0
	}
	return float64(matchCount) / float64(max(len(queryTerms), len(titleTerms)))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
