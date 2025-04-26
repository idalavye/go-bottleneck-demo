package api

import (
	"encoding/json"
	"fmt"
	"github.com/idagdelen/go-bottlenecks/internal/search"
	"net/http"
	"strconv"
	"time"
)

// Helper function to parse count parameter
func ParseCountParam(r *http.Request, defaultCount int) int {
	count := defaultCount
	countParam := r.URL.Query().Get("count")
	if countParam != "" {
		fmt.Sscanf(countParam, "%d", &count)
	}
	return count
}

// HandleSequential handles sequential processing requests
// @Summary Sequential Processing Demo
// @Description Demonstrates sequential processing
// @Tags performance
// @Param count query int false "Number of items to process"
// @Produce json
// @Success 200 {object} Response
// @Router /api/sequential [get]
func HandleSequential(w http.ResponseWriter, r *http.Request) {
	count := ParseCountParam(r, 10)

	resp := Response{
		Success: true,
		Message: "Sequential processing example - ready to run",
		Data: map[string]interface{}{
			"count":  count,
			"status": "prepared",
			"note":   "This endpoint is ready but actual processing is disabled",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleConcurrent handles concurrent processing requests
// @Summary Concurrent Processing Demo
// @Description Demonstrates concurrent processing with goroutines
// @Tags performance
// @Param count query int false "Number of items to process"
// @Produce json
// @Success 200 {object} Response
// @Router /api/concurrent [get]
func HandleConcurrent(w http.ResponseWriter, r *http.Request) {
	count := ParseCountParam(r, 10)

	resp := Response{
		Success: true,
		Message: "Concurrent processing example - ready to run",
		Data: map[string]interface{}{
			"count":  count,
			"status": "prepared",
			"note":   "This endpoint is ready but actual processing is disabled",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandlePool handles worker pool requests
// @Summary Worker Pool Demo
// @Description Demonstrates concurrent processing with a worker pool
// @Tags performance
// @Param count query int false "Number of items to process"
// @Param workers query int false "Number of workers"
// @Produce json
// @Success 200 {object} Response
// @Router /api/pool [get]
func HandlePool(w http.ResponseWriter, r *http.Request) {
	count := ParseCountParam(r, 10)

	workers := 4 // Default worker count
	workersParam := r.URL.Query().Get("workers")
	if workersParam != "" {
		fmt.Sscanf(workersParam, "%d", &workers)
	}

	resp := Response{
		Success: true,
		Message: "Worker pool example - ready to run",
		Data: map[string]interface{}{
			"count":   count,
			"workers": workers,
			"status":  "prepared",
			"note":    "This endpoint is ready but actual processing is disabled",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleSearch handles search
// @Summary Search Demo
// @Description Searching to keyword on vectorized db then populate it external services
// @Tags search
// @Router /api/search [get]
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("term")
	itemCount := r.URL.Query().Get("itemCount")

	parsedItemCount, err := strconv.Atoi(itemCount)
	if err != nil {
		return
	}

	products := search.SearchProducts(searchTerm, parsedItemCount)

	resp := Response{
		Success: true,
		Message: "Worker pool example - ready to run",
		Data: map[string]interface{}{
			"result": products,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleLeak handles goroutine leak demonstration requests
// @Summary Goroutine Leak Demo
// @Description Demonstrates a goroutine leak (use with caution)
// @Tags performance
// @Param count query int false "Number of items to process"
// @Param leak query bool false "Should create a leak"
// @Produce json
// @Success 200 {object} Response
// @Router /api/leak [get]
func HandleLeak(w http.ResponseWriter, r *http.Request) {
	count := ParseCountParam(r, 10)

	shouldLeak := false
	leakParam := r.URL.Query().Get("leak")
	if leakParam == "true" {
		shouldLeak = true
	}

	resp := Response{
		Success: true,
		Message: "Goroutine leak example - ready to run",
		Data: map[string]interface{}{
			"count":      count,
			"shouldLeak": shouldLeak,
			"status":     "prepared",
			"note":       "This endpoint is ready but actual processing is disabled",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleHealthCheck handles health check requests
// @Summary Health Check
// @Description Returns API health status
// @Tags system
// @Produce json
// @Success 200 {object} Response
// @Router /api/health [get]
func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Success: true,
		Message: "API is healthy",
		Data: map[string]interface{}{
			"status": "up",
			"time":   time.Now().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
