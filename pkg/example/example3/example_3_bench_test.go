package example3

import (
	"encoding/json"
	"testing"
)

var exampleStruct = Product{
	ID:     1,
	Name:   "Test User",
	Email:  "test@example.com",
	Tags:   []string{"go", "pool", "json", "benchmark"},
	Active: true,
	Score:  99.9,
	Meta:   map[string]interface{}{"role": "admin", "level": 5},
}

var exampleJSON, _ = json.Marshal(exampleStruct)

func BenchmarkParseJSONWithoutPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = parseJSONWithoutPool(exampleJSON)
	}
}

func BenchmarkParseJSONWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = parseJSONWithPool(exampleJSON)
	}
}
