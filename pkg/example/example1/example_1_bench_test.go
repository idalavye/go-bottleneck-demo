package example1

import (
	"math/rand"
	"testing"
)

// generateMatrix generates a 2D slice with given dimensions and random values
func generateMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(100)
		}
	}
	return matrix
}

func BenchmarkRowBasedSum(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RowBasedSum(matrix)
	}
}

func BenchmarkColumnBasedSum(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ColumnBasedSum(matrix)
	}
}
