package utils

import "strings"

// returns the content around the query, 100 characters before and after
// and adds ellipsis if needed
// the content is split by whole words, so the result is not exactly 200 characters
// but will not cut words in half
func GetContentAroundMatch(content, query string) string {
	// Find the first occurrence of the query in the content
	index := strings.Index(content, query)
	if index == -1 {
		return ""
	}

	start := index - 100
	if start < 0 {
		start = 0
	}
	end := index + 100
	if end > len(content) {
		end = len(content)
	}
	
	// extract the content around the query, without cutting words in half
	// find the first whole word before the start
	for start > 0 && content[start] != ' ' {
		start--
	}
	if start < 0 {
		start = 0
	}

	// find the first whole word after the end
	for end < len(content) && content[end] != ' ' {
		end++
	}
	if end > len(content) {
		end = len(content)
	}

	// Extract the content around the query and add ellipsis if needed
	res := content[start:end]
	if start > 0 {
		res = "..." + res
	}
	if end < len(content) {
		res = res + "..."
	}

	return res
}