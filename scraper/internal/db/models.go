package db

import "time"

type Page struct {
    ID        uint
    Title     string
    URL       string
    Language  string
    Content   string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewPage creates a new Page with default values
func NewPage() *Page {
    now := time.Now()
    return &Page{
        Language:  "en",
        CreatedAt: now,
        UpdatedAt: now,
    }
}

type QueryResult struct {
    ID    uint
    Query string
    Count int
}