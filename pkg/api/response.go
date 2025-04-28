package api

import "time"

// Response holds the unified API response structure
type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Duration  string      `json:"duration,omitempty"`
	StartTime time.Time   `json:"startTime,omitempty"`
	EndTime   time.Time   `json:"endTime,omitempty"`
}

// EnrichedProduct holds both the search score and detailed product info
// This struct is used to return all relevant info in the response
// (score + product details)
type EnrichedProduct struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       string  `json:"price"` // Price formatted as string with currency
	Score       float64 `json:"score"`
	Stock       int     `json:"stock"` // Stock quantity for the product
}
