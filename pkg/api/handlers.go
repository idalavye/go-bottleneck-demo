package api

import (
	"encoding/json"
	"net/http"
	"runtime/trace"
	"strconv"
	"time"

	"github.com/idagdelen/go-bottlenecks/internal/search"
)

// HandleSearch handles search
// @Summary Search Demo
// @Description Searches the vectorized database with the given keyword and enriches the results with external services
// @Tags search
// @Param term query string true "Keyword to search for"
// @Param itemCount query int false "Number of products to return (default: 10)"
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/search [get]
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	ctx, task := trace.NewTask(r.Context(), "HandleSearch")
	defer task.End()
	trace.WithRegion(ctx, "HandleSearch-Request", func() {
		searchTerm := r.URL.Query().Get("term")
		itemCount := r.URL.Query().Get("itemCount")

		parsedItemCount, err := strconv.Atoi(itemCount)
		if err != nil {
			parsedItemCount = 10 // Default value
		}

		products, totalSum := search.SearchProducts(searchTerm, parsedItemCount)

		resp := Response{
			Success: true,
			Message: "Search completed successfully",
			Data: map[string]interface{}{
				"result":   products,
				"totalSum": totalSum,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
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
