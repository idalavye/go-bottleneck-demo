package search

import (
	"container/heap"
	"math"
	"runtime"
	"sort"
)

/*
This function measures how similar the directions of two vectors are.
It is widely used in recommendation systems, search algorithms, and text similarity applications.
*/
func cosineSimilarity(a, b []float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Min-heap implementation for ScoredProduct
// Only keeps the top N scored products in memory

type scoredProductMinHeap []ScoredProduct

func (h scoredProductMinHeap) Len() int           { return len(h) }
func (h scoredProductMinHeap) Less(i, j int) bool { return h[i].Score < h[j].Score } // min-heap
func (h scoredProductMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *scoredProductMinHeap) Push(x interface{}) {
	*h = append(*h, x.(ScoredProduct))
}
func (h *scoredProductMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// SearchProducts optimized: only keeps top N results in memory using a min-heap
func SearchProducts(text string, pageSize int) ([]ScoredProduct, float64) {
	queryVector := getEmbedding(text)

	var scoredProducts []ScoredProduct
	for i, vec := range productVectors {
		score := cosineSimilarity(queryVector, vec)
		scoredProducts = append(scoredProducts, ScoredProduct{
			Product: productMetadata[i],
			Score:   score,
		})
	}

	sort.Slice(scoredProducts, func(i, j int) bool {
		return scoredProducts[i].Score > scoredProducts[j].Score
	})

	if pageSize < len(scoredProducts) {
		scoredProducts = scoredProducts[:pageSize]
	}

	return scoredProducts, SumRowMajor(productVectors)
}

// SearchProductsHeapOptimized: only keeps top N results in memory using a min-heap
// This method is memory efficient for large product sets.
/*
"Heap optimizasyonundan önce arama fonksiyonu belleğin %94'ünü kullanıyordu. Optimizasyon sonrası ise neredeyse hiç bellek kullanmıyor. Artık bottleneck başka bir noktada."

go tool pprof http://localhost:6060/debug/pprof/heap
*/
func SearchProductsHeapOptimized(text string, pageSize int) ([]ScoredProduct, float64) {
	queryVector := getEmbedding(text)

	h := &scoredProductMinHeap{}
	heap.Init(h)

	for i, vec := range productVectors {
		score := cosineSimilarity(queryVector, vec)
		item := ScoredProduct{
			Product: productMetadata[i],
			Score:   score,
		}
		if h.Len() < pageSize {
			heap.Push(h, item)
		} else if h.Len() > 0 && score > (*h)[0].Score {
			// Replace the smallest if current score is higher
			heap.Pop(h)
			heap.Push(h, item)
		}
	}

	// Extract results from heap and sort descending
	topResults := make([]ScoredProduct, h.Len())
	for i := len(topResults) - 1; i >= 0; i-- {
		topResults[i] = heap.Pop(h).(ScoredProduct)
	}

	return topResults, 0
}

// SearchProductsHeapOptimizedParallel: Calculates scores in parallel and keeps top N results using a min-heap
// This method is CPU efficient for large product sets.
/*
go tool pprof http://localhost:6060/debug/pprof/profile\?seconds\=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
*/
func SearchProductsHeapOptimizedParallel(text string, pageSize int) ([]ScoredProduct, float64) {
	queryVector := getEmbedding(text)

	numWorkers := runtime.NumCPU()
	jobs := make(chan int, len(productVectors))
	results := make(chan ScoredProduct, len(productVectors))

	// Worker goroutines for parallel score calculation
	for w := 0; w < numWorkers; w++ {
		go func() {
			for i := range jobs {
				score := cosineSimilarity(queryVector, productVectors[i])
				results <- ScoredProduct{
					Product: productMetadata[i],
					Score:   score,
				}
			}
		}()
	}

	for i := range productVectors {
		jobs <- i
	}
	close(jobs)

	allResults := make([]ScoredProduct, 0, len(productVectors))
	for i := 0; i < len(productVectors); i++ {
		allResults = append(allResults, <-results)
	}
	close(results)

	// Use a min-heap to keep only the top N results
	h := &scoredProductMinHeap{}
	heap.Init(h)
	for _, item := range allResults {
		if h.Len() < pageSize {
			heap.Push(h, item)
		} else if h.Len() > 0 && item.Score > (*h)[0].Score {
			// Replace the smallest if current score is higher
			heap.Pop(h)
			heap.Push(h, item)
		}
	}

	// Extract results from heap and sort descending
	topResults := make([]ScoredProduct, h.Len())
	for i := len(topResults) - 1; i >= 0; i-- {
		topResults[i] = heap.Pop(h).(ScoredProduct)
	}

	return topResults, 0
}

// Row-major traversal: Sums all elements row by row
func SumRowMajor(matrix [][]float64) float64 {
	var sum float64
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			sum += matrix[i][j]
		}
	}
	return sum
}

// Column-major traversal: Sums all elements column by column
func SumColMajor(matrix [][]float64) float64 {
	var sum float64
	if len(matrix) == 0 {
		return 0
	}
	for j := 0; j < len(matrix[0]); j++ {
		for i := 0; i < len(matrix); i++ {
			sum += matrix[i][j]
		}
	}
	return sum
}

// Parallel row-major traversal: Each row is summed in a separate goroutine, results are merged in main goroutine
func ParallelRowSumWithMerge(matrix [][]float64) float64 {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	rows := len(matrix)
	resultCh := make(chan float64, rows)

	// Her satır için ayrı bir goroutine açılıyor
	for i := 0; i < rows; i++ {
		go func(row []float64) {
			var rowSum float64
			for _, v := range row {
				rowSum += v
			}
			resultCh <- rowSum
		}(matrix[i])
	}

	var sum float64
	for i := 0; i < rows; i++ {
		sum += <-resultCh
	}
	close(resultCh)
	return sum
}

// Parallel row-major traversal: Sums all rows using a worker pool (worker sayısı = CPU sayısı)
func ParallelRowSumWithWorkers(matrix [][]float64) float64 {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	rows := len(matrix)
	numWorkers := runtime.NumCPU()
	rowCh := make(chan []float64, rows)
	resultCh := make(chan float64, rows)

	for w := 0; w < numWorkers; w++ {
		go func() {
			for row := range rowCh {
				var rowSum float64
				for _, v := range row {
					rowSum += v
				}
				resultCh <- rowSum
			}
		}()
	}

	for i := 0; i < rows; i++ {
		rowCh <- matrix[i]
	}
	close(rowCh)

	var sum float64
	for i := 0; i < rows; i++ {
		sum += <-resultCh
	}
	close(resultCh)
	return sum
}

// Row sums with escape to heap demonstration
// Returns a slice of row sums, which will escape to heap
func RowSumsEscapeHeap(matrix [][]float64) []float64 {
	rowSums := make([]float64, len(matrix))
	for i := 0; i < len(matrix); i++ {
		var sum float64
		for j := 0; j < len(matrix[i]); j++ {
			sum += matrix[i][j]
		}
		rowSums[i] = sum
	}
	return rowSums // rowSums slice backing array escapes to heap
}

// Sum using RowSumsEscapeHeap: sums all row sums
/*
go build -gcflags="-m" ./internal/search
*/
func SumWithRowSumsEscapeHeap(matrix [][]float64) float64 {
	rowSums := RowSumsEscapeHeap(matrix)
	var total float64
	for _, v := range rowSums {
		total += v
	}
	return total
}
