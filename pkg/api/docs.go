package api

// APIDocumentation contains the API documentation for the application
var APIDocumentation = map[string]interface{}{
	"info": map[string]string{
		"title":       "Go Bottlenecks API",
		"description": "API for demonstrating Go performance bottlenecks",
		"version":     "1.0.0",
	},
	"paths": map[string]interface{}{
		"/": map[string]interface{}{
			"get": map[string]string{
				"summary":     "Root endpoint",
				"description": "Returns a welcome message",
			},
		},
		"/api/docs": map[string]interface{}{
			"get": map[string]string{
				"summary":     "API Documentation",
				"description": "Returns API documentation in JSON format",
			},
		},
		"/api/sequential": map[string]interface{}{
			"get": map[string]string{
				"summary":     "Sequential Processing Demo",
				"description": "Demonstrates sequential processing",
				"parameters":  "?count=10 (optional, number of items to process)",
			},
		},
		"/api/concurrent": map[string]interface{}{
			"get": map[string]string{
				"summary":     "Concurrent Processing Demo",
				"description": "Demonstrates concurrent processing with goroutines",
				"parameters":  "?count=10 (optional, number of items to process)",
			},
		},
		"/api/pool": map[string]interface{}{
			"get": map[string]string{
				"summary":     "Worker Pool Demo",
				"description": "Demonstrates concurrent processing with a worker pool",
				"parameters":  "?count=10&workers=4 (optional parameters)",
			},
		},
		"/api/leak": map[string]interface{}{
			"get": map[string]string{
				"summary":     "Goroutine Leak Demo",
				"description": "Demonstrates a goroutine leak (use with caution)",
				"parameters":  "?count=100&leak=false (optional parameters)",
			},
		},
	},
}
