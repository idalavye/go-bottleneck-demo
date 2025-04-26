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
