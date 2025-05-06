package api

import (
	"testing"

	"github.com/idagdelen/go-bottlenecks/internal/search"
)

/*
go test -bench=BenchmarkEnrichProductsWithDetailsAndAdWorkerPool -run=^$ ./pkg/api -memprofile=mem.prof -trace=trace.out -cpu 10 -benchtime 3s
go tool pprof -http=:8080 ./cpu.prof
go tool pprof -http=:8080 ./mem.prof
go tool trace trace.out
*/
func BenchmarkEnrichProductsWithDetailsAndAdWorkerPool(b *testing.B) {
	// 100 ürünlük sahte bir ürün listesi oluştur
	products := make([]search.ScoredProduct, 100)
	for i := 0; i < 100; i++ {
		products[i] = search.ScoredProduct{
			Product: &search.Product{ID: i + 1, Name: "Product"},
			Score:   float64(i) * 1.1,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enrichProductsWithDetailsAndAdWorkerPool(products)
	}
}

/*
go test -bench=BenchmarkEnrichProductsWithDetailsAndAd -run=^$ ./pkg/api -memprofile=mem.prof -trace=trace.out -cpu 10 -benchtime 3s
go tool pprof -http=:8080 ./cpu.prof
go tool pprof -http=:8080 ./mem.prof
go tool trace trace.out
*/
func BenchmarkEnrichProductsWithDetailsAndAd(b *testing.B) {
	// 100 ürünlük sahte bir ürün listesi oluştur
	products := make([]search.ScoredProduct, 100)
	for i := 0; i < 100; i++ {
		products[i] = search.ScoredProduct{
			Product: &search.Product{ID: i + 1, Name: "Product"},
			Score:   float64(i) * 1.1,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enrichProductsWithDetailsAndAd(products)
	}
}
