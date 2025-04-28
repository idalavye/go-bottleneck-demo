package search

import (
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

// Part 1
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
