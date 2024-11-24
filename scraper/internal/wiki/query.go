package wiki

import (
    "fmt"
    "net/url"
    "regexp"
    "strings"
)

// preprocessQuery processes the input query string by trimming whitespace,
// and extracting key terms if it is a programming-related question.
// It returns the processed query string and 
// a boolean indicating whether the query is valid.
//
// Parameters:
//   - query: The input query string to be processed.
//
// Returns:
//   - string: The processed query string or an empty string if the query is invalid.
//   - bool: A boolean indicating whether the query is valid (true) or not (false).
func preprocessQuery(query string) (string, bool) {
    query = strings.TrimSpace(query)
    if query == "" {
        return "", false
    }

    if isProgrammingQuestion(query) {
        terms := extractKeyTerms(query)
        if len(terms) > 0 {
            return terms[len(terms)-1], true
        }
        return "", false
    }

    return cleanQuery(query), true
}

// isProgrammingQuestion checks if a given query string matches any of the predefined
// patterns that indicate it is a programming-related question. The function uses
// regular expressions to perform case-insensitive matching against a set of common
// programming-related queries.
//
// Parameters:
//   - query: A string representing the user's query.
//
// Returns:
//   - bool: Returns true if the query matches any of the programming-related patterns,
//           otherwise returns false.
func isProgrammingQuestion(query string) bool {
    patterns := []string{
        `(?i)what is the best programming language`,
        `(?i)how to.*(?:code|program|develop|implement|write|use|run)`,
        `(?i)what are the.*(?:features|capabilities|requirements)`,
    }

    for _, pattern := range patterns {
        if matched, _ := regexp.MatchString(pattern, query); matched {
            return true
        }
    }
    return false
}

// extractKeyTerms processes a given query string to extract key terms.
// It converts the query to lowercase, removes specific prefixes, splits the query into words,
// and filters out stop words.
//
// Parameters:
//   - query: A string representing the query to be processed.
//
// Returns:
//   - A slice of strings containing the key terms extracted from the query.
func extractKeyTerms(query string) []string {
    query = strings.ToLower(query)
    prefixes := []string{
        "what is the best programming language for",
        "how to program", "how to implement",
        "what are the features of",
        "how do i use",
    }
    
    for _, prefix := range prefixes {
        query = strings.TrimPrefix(query, prefix)
    }
    
    words := strings.Fields(query)
    var terms []string
    
    for _, word := range words {
        if !isStopWord(word) {
            terms = append(terms, word)
        }
    }
    return terms
}

// cleanQuery takes a search query string and removes common prefixes and suffixes
// to clean up the query. It converts the query to lowercase, trims leading and
// trailing spaces, and removes specified prefixes and suffixes. Finally, it formats
// the cleaned query for use in a URL.
//
// Parameters:
//   - query: The original search query string.
//
// Returns:
//   - A cleaned and formatted query string suitable for use in a URL.
func cleanQuery(query string) string {
    prefixes := []string{
        "what is", "what are", "how to", "is ", "are ", "why ",
        "when ", "where ", "who ", "which ", "how ", "can ",
    }
    
    query = strings.ToLower(strings.TrimSpace(query))
    for _, prefix := range prefixes {
        if strings.HasPrefix(query, prefix) {
            query = strings.TrimSpace(strings.TrimPrefix(query, prefix))
        }
    }

    suffixes := []string{
        " in 2024", " 2024", " 2025",
        " stock", " wiki", " wikipedia",
    }
    for _, suffix := range suffixes {
        if strings.HasSuffix(query, suffix) {
            query = strings.TrimSpace(strings.TrimSuffix(query, suffix))
        }
    }

    return formatQueryForURL(query)
}

// formatQueryForURL formats a query string for use in a URL.
// It splits the query into words, capitalizes the first letter of each word,
// and joins them with underscores. Stop words are not capitalized unless they
// are the first word in the query.
//
// Parameters:
//   query - The query string to format.
//
// Returns:
//   A string formatted for use in a URL.
func formatQueryForURL(query string) string {
    words := strings.Fields(query)
    for i, word := range words {
        if i == 0 || !isStopWord(word) {
            words[i] = strings.Title(word)
        }
    }
    return strings.Join(words, "_")
}

// isStopWord checks if a given word is a stop word.
// Stop words are common words that are typically filtered out in text processing.
// The function returns true if the word is a stop word, and false otherwise.
// The comparison is case-insensitive.
//
// Parameters:
//   - word: The word to check.
//
// Returns:
//   - bool: True if the word is a stop word, false otherwise.
func isStopWord(word string) bool {
    stopWords := map[string]bool{
        "a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
        "be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
        "in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
        "that": true, "the": true, "to": true, "was": true, "were": true,
        "will": true, "with": true,
    }
    return stopWords[strings.ToLower(word)]
}

// buildSearchURLs constructs a list of Wikipedia search URLs based on the provided queries.
// It returns nil if the processedQuery is empty.
//
// Parameters:
//   - processedQuery: A string representing the processed version of the search query.
//   - originalQuery: A string representing the original search query.
//
// Returns:
//   A slice of strings containing URLs for searching Wikipedia using the provided queries.
func buildSearchURLs(processedQuery, originalQuery string) []string {
    if processedQuery == "" {
        return nil
    }

    return []string{
        // Direct search using Wikipedia's search
        fmt.Sprintf("https://en.wikipedia.org/w/index.php?search=%s&title=Special:Search", 
            url.QueryEscape(originalQuery)),
        // Try direct article lookup
        "https://en.wikipedia.org/wiki/" + url.QueryEscape(processedQuery),
        // Try original query as fallback
        "https://en.wikipedia.org/wiki/" + url.QueryEscape(strings.ReplaceAll(originalQuery, " ", "_")),
    }
}