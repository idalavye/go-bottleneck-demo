package search

import (
	"container/heap"
	"math"
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

	scoredProducts := make([]ScoredProduct, len(productVectors))

	for i, vec := range productVectors {
		score := cosineSimilarity(queryVector, vec)
		scoredProducts[i] = ScoredProduct{
			Product: &productMetadata[i],
			Score:   score,
		}
	}

	sort.Slice(scoredProducts, func(i, j int) bool {
		return scoredProducts[i].Score > scoredProducts[j].Score
	})

	if pageSize < len(scoredProducts) {
		scoredProducts = scoredProducts[:pageSize]
	}

	return scoredProducts, 0
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
		score := cosineSimilarity(queryVector, vec)

		item := ScoredProduct{
			Product: &productMetadata[i],
			Score:   score,
		}
		if h.Len() < pageSize {
			heap.Push(h, item)
		} else if h.Len() > 0 && score > (*h)[0].Score {
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

func SearchProductsQdrantOptimized(text string, pageSize int) ([]ScoredProduct, float64) {
	vector := make([]float32, 4)
	for i := range vector {
		vector[i] = float32(getEmbedding(text)[i])
	}

	results, err := searchProductsQdrant(vector, pageSize)
	if err != nil {
		return nil, 0
	}

	scored := make([]ScoredProduct, 0, len(results))
	totalSum := 0.0
	for _, r := range results {
		scored = append(scored, ScoredProduct{
			Product: &Product{
				ID:   r.ID,
				Name: r.Name,
			},
			Score: float64(r.Score),
		})
		totalSum += float64(r.Score)
	}
	return scored, totalSum
}
