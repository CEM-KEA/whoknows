package utils

import "strings"

// returns the content around the query, 100 characters before and after
func GetContentAroundMatch(content, query string) string {
	// Find the first occurrence of the query in the content
	index := strings.Index(content, query)
	if index == -1 {
		return ""
	}

	// Get the start and end indexes of the content around the query
	start := index - 100
	if start < 0 {
		start = 0
	}
	end := index + 100
	if end > len(content) {
		end = len(content)
	}

	// Get the content around the query
	return content[start:end]
}