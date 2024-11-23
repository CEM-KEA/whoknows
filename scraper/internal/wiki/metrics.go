package wiki

import (
    "log"
    "time"
)

type Metrics struct {
    QueriesProcessed  int
    SuccessfulScrapes int
    FailedScrapes     int
    StartTime         time.Time
    EndTime           time.Time
}

func NewMetrics() *Metrics {
    return &Metrics{
        StartTime: time.Now(),
    }
}

func (m *Metrics) Log() {
    m.EndTime = time.Now()
    duration := m.EndTime.Sub(m.StartTime)
    log.Printf("Scraper finished in %v", duration)
    log.Printf("Queries processed: %d", m.QueriesProcessed)
    log.Printf("Successful scrapes: %d", m.SuccessfulScrapes)
    log.Printf("Failed scrapes: %d", m.FailedScrapes)
}