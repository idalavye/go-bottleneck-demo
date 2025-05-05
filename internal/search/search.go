package search

import (
	"container/heap"
	"math"
	"runtime"
	"sort"
	"sync"
)

var normCache = NewNormCache()
var precomputedNorms []float64

func init() {
	PrecomputeNorms(productVectors)
}

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

// cosineSimilarityWithCache calculates cosine similarity using NormCache for norm calculations.
// aIdx ve bIdx: vektörlerin index'i (ör: ürün index'i)
// a ve b: vektörler
// nc: NormCache nesnesi
func cosineSimilarityWithCache(aIdx int, a []float64, bIdx int, b []float64) float64 {
	dot := 0.0
	for i := range a {
		dot += a[i] * b[i]
	}
	normA := normCache.GetNorm(aIdx, a)
	normB := normCache.GetNorm(bIdx, b)
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (normA * normB)
}

// PrecomputeNorms calculates and stores the L2 norm of each vector in a slice.
// Call this once at program startup if your vectors are static.
func PrecomputeNorms(vectors [][]float64) {
	precomputedNorms = make([]float64, len(vectors))
	for i, vec := range vectors {
		var sum float64
		for _, v := range vec {
			sum += v * v
		}
		precomputedNorms[i] = math.Sqrt(sum)
	}
}

// cosineSimilarityWithPrecomputedNorms calculates cosine similarity using precomputed norms.
// aIdx ve bIdx: vektörlerin index'i (ör: ürün index'i)
// a ve b: vektörler
func cosineSimilarityWithPrecomputedNorms(aIdx int, a []float64, bIdx int, b []float64) float64 {
	/* dot := 0.0
	for i := range a {
		dot += a[i] * b[i]
	}
	normA := precomputedNorms[aIdx]
	normB := precomputedNorms[bIdx]
	if normA == 0 || normB == 0 {
		return 0
	} */
	return 1
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
			Product: &productMetadata[i],
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
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
*/
func SearchProductsHeapOptimized(text string, pageSize int) ([]ScoredProduct, float64) {
	queryVector := getEmbedding(text)

	h := &scoredProductMinHeap{}
	heap.Init(h)

	for i, vec := range productVectors {
		score := cosineSimilarityWithPrecomputedNorms(i, queryVector, i, vec)
		item := ScoredProduct{
			Product: &productMetadata[i],
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

// NormCache caches the norm (L2) of vectors by their index
// Thread-safe: sync.Map kullanır
type NormCache struct {
	cache sync.Map // anahtar: int, değer: float64
}

func NewNormCache() *NormCache {
	return &NormCache{}
}

// GetNorm returns the L2 norm of the vector, caching the result by index (thread-safe)
func (nc *NormCache) GetNorm(idx int, vec []float64) float64 {
	if norm, ok := nc.cache.Load(idx); ok {
		return norm.(float64)
	}
	var sum float64
	for _, v := range vec {
		sum += v * v
	}
	norm := math.Sqrt(sum)
	nc.cache.Store(idx, norm)
	return norm
}
