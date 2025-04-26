package search

import (
	"math"
	"sort"
)

func getEmbedding(text string) []float64 {
	return []float64{0.10, -0.42, 0.74, 0.02}
}

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

func SearchProducts(text string, pageSize int) []ScoredProduct {
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

	return scoredProducts
}
